.PHONY: default help clone fetch sync checkout init-repo-envs envs doctor ssh-helper ssh status up down restart hotreload ps logs

default: help

help:
	@cd git-controller && go run . help

clone fetch sync checkout init-repo-envs doctor ssh-helper status:
	@cd git-controller && go run . $@

envs: init-repo-envs
ssh: ssh-helper

up down restart hotreload ps logs:
	@$(MAKE) -C local-orchestrator $@
