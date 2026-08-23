.DEFAULT_GOAL := help

# ── Compose file paths ─────────────────────────────────────────
COMPOSE_FILE   := docker-compose.yml
HOTRELOAD_FILE := hotreload/docker-compose.yml

# ── Colours (ANSI) ─────────────────────────────────────────────
CYAN  := \033[36m
GREEN := \033[32m
BOLD  := \033[1m
RESET := \033[0m

# ═══════════════════════════════════════════════════════════════
#  Stack Lifecycle
# ═══════════════════════════════════════════════════════════════

.PHONY: up
up: ## Start all services in the background (standard mode)
	@printf "$(CYAN)▶ Starting local stack …$(RESET)\n"
	@docker compose -f $(COMPOSE_FILE) up -d --remove-orphans

.PHONY: down
down: ## Stop all services
	@printf "$(CYAN)▶ Stopping local stack …$(RESET)\n"
	@docker compose -f $(COMPOSE_FILE) down

.PHONY: down-v
down-v: ## Stop all services and remove volumes
	@printf "$(CYAN)▶ Stopping local stack and removing volumes …$(RESET)\n"
	@docker compose -f $(COMPOSE_FILE) down -v --remove-orphans
	@docker compose -f $(HOTRELOAD_FILE) down -v --remove-orphans 2>/dev/null || true

.PHONY: restart
restart: down up ## Restart all services

.PHONY: hotreload up-hotreload dev
hotreload up-hotreload dev: ## Start all services with live source mounting & hot reload
	@printf "$(CYAN)▶ Starting local stack in hotreload mode …$(RESET)\n"
	@docker compose -f $(HOTRELOAD_FILE) up -d --remove-orphans

.PHONY: down-hotreload
down-hotreload: ## Stop hotreload stack
	@printf "$(CYAN)▶ Stopping hotreload stack …$(RESET)\n"
	@docker compose -f $(HOTRELOAD_FILE) down

.PHONY: ps
ps: ## Show status of running containers
	@docker compose -f $(COMPOSE_FILE) ps

.PHONY: logs
logs: ## Tail logs for all services (Ctrl-C to stop)
	@docker compose -f $(COMPOSE_FILE) logs -f --tail=50

# ═══════════════════════════════════════════════════════════════
#  Single Service (e.g. make service-up S=placeholder1-service)
# ═══════════════════════════════════════════════════════════════

.PHONY: service-up
service-up: ## Start a single service (S=<name>)
	@test -n "$(S)" || (printf "Usage: make service-up S=<service-name>\n" && exit 1)
	@printf "$(CYAN)▶ Starting $(S) …$(RESET)\n"
	@docker compose -f $(COMPOSE_FILE) up -d --remove-orphans $(S)

.PHONY: service-down
service-down: ## Stop a single service (S=<name>)
	@test -n "$(S)" || (printf "Usage: make service-down S=<service-name>\n" && exit 1)
	@docker compose -f $(COMPOSE_FILE) stop $(S)

.PHONY: service-logs
service-logs: ## Tail logs for a single service (S=<name>)
	@test -n "$(S)" || (printf "Usage: make service-logs S=<service-name>\n" && exit 1)
	@docker compose -f $(COMPOSE_FILE) logs -f --tail=50 $(S)

# ═══════════════════════════════════════════════════════════════
#  Utilities & Maintenance
# ═══════════════════════════════════════════════════════════════

.PHONY: validate
validate: ## Validate compose configurations (standard and hotreload)
	@printf "$(CYAN)Validating compose configurations …$(RESET)\n"
	@docker compose -f $(COMPOSE_FILE) config -q && printf "  $(GREEN)✔$(RESET) $(COMPOSE_FILE) valid\n"
	@docker compose -f $(HOTRELOAD_FILE) config -q && printf "  $(GREEN)✔$(RESET) $(HOTRELOAD_FILE) valid\n"
	@printf "$(GREEN)All compose files are valid.$(RESET)\n"

.PHONY: clean
clean: down ## Stop everything and remove volumes and orphaned containers
	@docker compose -f $(COMPOSE_FILE) down --volumes --remove-orphans 2>/dev/null || true
	@docker compose -f $(HOTRELOAD_FILE) down --volumes --remove-orphans 2>/dev/null || true
	@printf "$(GREEN)✔ Cleaned up.$(RESET)\n"

.PHONY: help
help: ## Show this help
	@printf "$(BOLD)Local Orchestrator — Available Targets$(RESET)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-22s$(RESET) %s\n", $$1, $$2}'
	@printf "\nExample: $(BOLD)make up$(RESET) or $(BOLD)make hotreload$(RESET) or $(BOLD)make service-up S=placeholder1-service$(RESET)\n"
