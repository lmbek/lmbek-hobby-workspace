# Workspace Controller Architecture

## Overview

This repository is structured as a command-based workspace controller for local development environments.

## Structure

### commands/
User-facing operations. Each folder represents a command that can be executed.

- start → shows execution plan for services and decoupled infrastructure
- sync → materializes system locally
- validate → checks system consistency

### internal/
Shared system logic used by all commands.

This folder is NOT a utility folder.

It is an architecture boundary that contains:
- system parsing
- runtime logic
- git orchestration logic
- validation rules

## Why "internal"?

This is a Go language convention.

It enforces that:
- shared logic cannot be imported from outside this repository
- commands must use internal APIs instead of duplicating logic

## Rule

If logic is reused by more than one command, it belongs in internal/.

---

### Command Interface

All interactions happen through the `workspace-controller` binary:

- `workspace-controller start`
- `workspace-controller sync`
- `workspace-controller validate`
- `workspace-controller help`

*Note: For a general overview and getting started guide, see [../README.md](../README.md). Always keep the README updated.*