"""Custom Ansible filter to flatten system-definition.yaml repos into a list."""


def flatten_repos(repos, workspace_root):
    """Flatten the repos map from system-definition.yaml into a list of dicts.

    Each category can be either:
      - A flat/singleton repo (has 'repository' key directly)
      - A map of named repos

    Returns a list of dicts with keys:
      category, name, display_name, repository, version, target_path
    """
    import os

    result = []
    for cat_name, cat_value in (repos or {}).items():
        if not cat_value:
            continue

        cat_dir = os.path.join(workspace_root, "git-repositories", cat_name)

        # Check if this is a flat/singleton category (has 'repository' key)
        if isinstance(cat_value, dict) and "repository" in cat_value:
            repo_url = cat_value.get("repository", "")
            version = cat_value.get("version", "main")
            if repo_url and "@company" not in repo_url:
                result.append({
                    "category": cat_name,
                    "name": "",
                    "display_name": cat_name,
                    "repository": repo_url,
                    "version": version,
                    "target_path": cat_dir,
                })
        else:
            # Map of named repos
            for name, comp in cat_value.items():
                if not isinstance(comp, dict):
                    continue
                repo_url = comp.get("repository", "")
                version = comp.get("version", "main")
                if repo_url and "@company" not in repo_url:
                    result.append({
                        "category": cat_name,
                        "name": name,
                        "display_name": name if name else cat_name,
                        "repository": repo_url,
                        "version": version,
                        "target_path": os.path.join(cat_dir, name),
                    })

    return result


class FilterModule:
    """Ansible filter plugin."""

    def filters(self):
        return {
            "flatten_repos": flatten_repos,
        }
