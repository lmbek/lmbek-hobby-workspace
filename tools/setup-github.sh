#!/usr/bin/env bash
set -euo pipefail

readonly GITHUB_OWNER="${GITHUB_OWNER:-lmbek}"
readonly PLATFORM_REPOSITORY="${GITHUB_OWNER}/lmbek-hobby-platform"
WORKSPACE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly WORKSPACE_ROOT
readonly TFVARS_PATH="${WORKSPACE_ROOT}/git-repositories/infrastructure/iac/terraform.tfvars"
readonly TFVARS_EXAMPLE_PATH="${TFVARS_PATH}.example"

readonly -a REPOSITORIES=(
  "lmbek-hobby-deployment"
  "lmbek-web-frontend-placeholder"
  "lmbek-hobby-placeholder1-service"
  "lmbek-hobby-placeholder2-service"
  "lmbek-hobby-infrastructure"
  "lmbek-hobby-platform"
  "lmbek-hobby-servers"
  "lmbek-hobby-observability"
  "lmbek-hobby-docs"
)

readonly -a IMAGE_REPOSITORIES=(
  "lmbek-web-frontend-placeholder:lmbek-hobby-web-frontend"
  "lmbek-hobby-placeholder1-service:lmbek-hobby-placeholder1-service"
  "lmbek-hobby-placeholder2-service:lmbek-hobby-placeholder2-service"
  "lmbek-hobby-docs:lmbek-hobby-docs"
)

readonly -a PROTECTED_REPOSITORIES=(
  "lmbek-web-frontend-placeholder:Test, build, scan, and publish"
  "lmbek-hobby-placeholder1-service:Test, build, scan, and publish"
  "lmbek-hobby-placeholder2-service:Test, build, scan, and publish"
  "lmbek-hobby-docs:Test, build, scan, and publish"
  "lmbek-hobby-infrastructure:terraform"
  "lmbek-hobby-platform:kustomize"
)

install_github_cli() {
  if command -v gh >/dev/null 2>&1; then
    return
  fi

  if ! command -v apt-get >/dev/null 2>&1; then
    printf 'GitHub CLI is missing and automatic installation currently supports Debian/Ubuntu only.\n' >&2
    exit 1
  fi

  local sudo_command=()
  if [[ "$(id -u)" -ne 0 ]]; then
    command -v sudo >/dev/null 2>&1 || {
      printf 'GitHub CLI installation requires root or sudo.\n' >&2
      exit 1
    }
    sudo_command=(sudo)
  fi

  printf '==> Installing GitHub CLI from the official repository...\n'
  "${sudo_command[@]}" apt-get update -y
  "${sudo_command[@]}" apt-get install -y ca-certificates curl
  "${sudo_command[@]}" install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg |
    "${sudo_command[@]}" tee /etc/apt/keyrings/githubcli-archive-keyring.gpg >/dev/null
  "${sudo_command[@]}" chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
  printf 'deb [arch=%s signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\n' \
    "$(dpkg --print-architecture)" |
    "${sudo_command[@]}" tee /etc/apt/sources.list.d/github-cli.list >/dev/null
  "${sudo_command[@]}" apt-get update -y
  "${sudo_command[@]}" apt-get install -y gh
}

authenticate_github_cli() {
  install_github_cli
  if gh auth status --hostname github.com >/dev/null 2>&1; then
    local scopes=""
    scopes=",$(gh api --include user 2>/dev/null |
      awk '!found && tolower($1) == "x-oauth-scopes:" { $1 = ""; sub(/^ /, ""); print; found = 1 }' |
      tr -d '\r '),"
    if [[ -z "${GH_TOKEN:-}" && ("$scopes" != *",repo,"* || "$scopes" != *",read:packages,"*) ]]; then
      printf '==> Expanding GitHub CLI authorization for repository and package administration...\n'
      gh auth refresh --hostname github.com --scopes repo,workflow,read:packages
    fi
    return
  fi

  if [[ -n "${GH_TOKEN:-}" ]]; then
    gh auth status --hostname github.com >/dev/null
    return
  fi

  printf '==> Authorize GitHub CLI in the browser. This is the only interactive setup step.\n'
  gh auth login --hostname github.com --git-protocol ssh --web \
    --scopes repo,workflow,read:packages
}

configure_actions() {
  local repository
  for repository in "${REPOSITORIES[@]}"; do
    printf '==> Enabling least-privilege Actions defaults for %s/%s...\n' "$GITHUB_OWNER" "$repository"
    gh api --method PUT "repos/${GITHUB_OWNER}/${repository}/actions/permissions" \
      --input - >/dev/null <<'JSON'
{"enabled":true,"allowed_actions":"all"}
JSON
    gh api --method PUT "repos/${GITHUB_OWNER}/${repository}/actions/permissions/workflow" \
      --input - >/dev/null <<'JSON'
{"default_workflow_permissions":"read","can_approve_pull_request_reviews":false}
JSON
  done
}

configure_production_environment() {
  printf '==> Restricting the platform production environment to main...\n'
  gh api --method PUT "repos/${PLATFORM_REPOSITORY}/environments/production" \
    --input - >/dev/null <<'JSON'
{"wait_timer":0,"reviewers":[],"deployment_branch_policy":{"protected_branches":false,"custom_branch_policies":true}}
JSON

  if ! gh api "repos/${PLATFORM_REPOSITORY}/environments/production/deployment-branch-policies" \
    --jq '.branch_policies[].name' | grep -Fxq main; then
    gh api --method POST "repos/${PLATFORM_REPOSITORY}/environments/production/deployment-branch-policies" \
      --field name=main --field type=branch >/dev/null
  fi
}

package_endpoint() {
  local package_name="$1"
  local owner_type
  owner_type="$(gh api "users/${GITHUB_OWNER}" --jq .type)"
  if [[ "$owner_type" == "Organization" ]]; then
    printf 'orgs/%s/packages/container/%s' "$GITHUB_OWNER" "$package_name"
  else
    printf 'user/packages/container/%s' "$package_name"
  fi
}

wait_for_dispatched_workflow() {
  local repository="$1"
  local previous_run_id="$2"
  local run_id=""
  local attempt=0

  while ((attempt < 30)); do
    attempt=$((attempt + 1))
    run_id="$(gh run list --repo "${GITHUB_OWNER}/${repository}" --workflow ci-cd.yml \
      --event workflow_dispatch --branch main --limit 1 --json databaseId --jq '.[0].databaseId // empty')"
    if [[ -n "$run_id" && "$run_id" != "$previous_run_id" ]]; then
      gh run watch "$run_id" --repo "${GITHUB_OWNER}/${repository}" --exit-status
      return
    fi
    sleep 2
  done

  printf 'Timed out waiting for the dispatched workflow in %s/%s.\n' "$GITHUB_OWNER" "$repository" >&2
  exit 1
}

configure_private_registry() {
  local pull_token="${GHCR_PULL_TOKEN:-}"
  if [[ -z "$pull_token" ]] && gh secret list --repo "$PLATFORM_REPOSITORY" --json name \
    --jq '.[].name' | grep -Fxq GHCR_PULL_TOKEN; then
    printf '==> Keeping the existing encrypted GHCR_PULL_TOKEN platform secret.\n'
    return
  fi
  if [[ -z "$pull_token" ]]; then
    if [[ ! -t 0 ]]; then
      printf 'Set GHCR_PULL_TOKEN for non-interactive setup.\n' >&2
      exit 1
    fi
    printf 'Paste a classic GitHub token with read:packages (input is hidden): ' >&2
    read -r -s pull_token || true
    printf '\n' >&2
  fi
  if [[ -z "$pull_token" ]]; then
    printf 'A GHCR pull token is required because the packages remain private.\n' >&2
    exit 1
  fi

  printf '==> Storing the private GHCR pull credential in the platform Actions secret...\n'
  printf '%s' "$pull_token" | gh secret set GHCR_PULL_TOKEN --repo "$PLATFORM_REPOSITORY"
  pull_token=""
}

publish_packages() {
  local mapping repository package_name endpoint previous_run_id
  for mapping in "${IMAGE_REPOSITORIES[@]}"; do
    repository="${mapping%%:*}"
    package_name="${mapping#*:}"
    endpoint="$(package_endpoint "$package_name")"

    if ! gh api "$endpoint" >/dev/null 2>&1; then
      printf '==> Running the initial image workflow for %s/%s...\n' "$GITHUB_OWNER" "$repository"
      previous_run_id="$(gh run list --repo "${GITHUB_OWNER}/${repository}" --workflow ci-cd.yml \
        --event workflow_dispatch --branch main --limit 1 --json databaseId --jq '.[0].databaseId // empty')"
      gh workflow run ci-cd.yml --repo "${GITHUB_OWNER}/${repository}" --ref main
      wait_for_dispatched_workflow "$repository" "$previous_run_id"
    fi

    printf '==> Verified private package ghcr.io/%s/%s.\n' "$GITHUB_OWNER" "$package_name"
  done
}

configure_ruleset() {
  local repository="$1"
  local required_check="$2"
  local ruleset_name="Protect main with CI"
  local ruleset_id endpoint bypass_actors

  ruleset_id="$(gh api "repos/${GITHUB_OWNER}/${repository}/rulesets" --paginate \
    --jq ".[] | select(.name == \"${ruleset_name}\") | .id" | head -n 1)"
  if [[ -n "$ruleset_id" ]]; then
    endpoint="repos/${GITHUB_OWNER}/${repository}/rulesets/${ruleset_id}"
  else
    endpoint="repos/${GITHUB_OWNER}/${repository}/rulesets"
  fi

  bypass_actors='[]'
  if [[ "$repository" == "lmbek-hobby-platform" ]]; then
    bypass_actors='[{"actor_id":15368,"actor_type":"Integration","bypass_mode":"always"}]'
  fi

  printf '==> Protecting main in %s/%s with required check "%s"...\n' \
    "$GITHUB_OWNER" "$repository" "$required_check"
  gh api --method "$(if [[ -n "$ruleset_id" ]]; then printf PUT; else printf POST; fi)" "$endpoint" \
    --input - >/dev/null <<JSON
{
  "name":"${ruleset_name}",
  "target":"branch",
  "enforcement":"active",
  "bypass_actors":${bypass_actors},
  "conditions":{"ref_name":{"include":["refs/heads/main"],"exclude":[]}},
  "rules":[
    {"type":"deletion"},
    {"type":"non_fast_forward"},
    {"type":"required_status_checks","parameters":{"strict_required_status_checks_policy":true,"required_status_checks":[{"context":"${required_check}"}]}}
  ]
}
JSON
}

configure_rulesets() {
  local mapping repository required_check
  for mapping in "${PROTECTED_REPOSITORIES[@]}"; do
    repository="${mapping%%:*}"
    required_check="${mapping#*:}"
    configure_ruleset "$repository" "$required_check"
  done
}

bootstrap() {
  authenticate_github_cli
  configure_actions
  configure_production_environment
  configure_private_registry
  publish_packages
  configure_rulesets
  printf '\nGitHub setup complete. The command is safe to rerun.\n'
}

write_runner_token() {
  authenticate_github_cli
  if [[ ! -f "$TFVARS_PATH" ]]; then
    cp "$TFVARS_EXAMPLE_PATH" "$TFVARS_PATH"
    printf 'Created %s from the example; fill its non-GitHub placeholders before Terraform apply.\n' "$TFVARS_PATH"
  fi

  local token temporary_file
  token="$(gh api --method POST "repos/${PLATFORM_REPOSITORY}/actions/runners/registration-token" --jq .token)"
  temporary_file="${TFVARS_PATH}.tmp.$$"
  trap 'rm -f "${temporary_file:-}"' EXIT
  awk -v token="$token" '
    /^[[:space:]]*github_runner_token[[:space:]]*=/ {
      print "github_runner_token = \"" token "\""
      replaced = 1
      next
    }
    { print }
    END {
      if (!replaced) {
        print ""
        print "github_runner_token = \"" token "\""
      }
    }
  ' "$TFVARS_PATH" > "$temporary_file"
  chmod 600 "$temporary_file"
  mv "$temporary_file" "$TFVARS_PATH"
  trap - EXIT
  printf 'Wrote a new one-hour, single-use runner token to the ignored terraform.tfvars. Apply Terraform now.\n'
}

case "${1:-bootstrap}" in
  bootstrap)
    bootstrap
    ;;
  runner-token)
    write_runner_token
    ;;
  *)
    printf 'Usage: %s {bootstrap|runner-token}\n' "$0" >&2
    exit 2
    ;;
esac