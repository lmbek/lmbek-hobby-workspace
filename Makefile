.PHONY: default help clone fetch sync checkout init-repo-envs envs doctor ssh-helper ssh status up down restart hotreload ps logs down-v

default: help

help:
	@cd git-controller && go run . help

clone fetch sync checkout init-repo-envs doctor ssh-helper status:
	@cd git-controller && go run . $@

envs: init-repo-envs
ssh: ssh-helper

up:
	@$(MAKE) -C local-orchestrator up

down:
	@$(MAKE) -C local-orchestrator down

restart:
	@$(MAKE) -C local-orchestrator restart

hotreload:
	@$(MAKE) -C local-orchestrator hotreload

ps:
	@$(MAKE) -C local-orchestrator ps

logs:
	@$(MAKE) -C local-orchestrator logs

down-v:
	@$(MAKE) -C local-orchestrator down-v