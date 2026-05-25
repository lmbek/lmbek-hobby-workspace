# System Controller Architecture

## Overview

This repository is structured as a command-based system controller for local development environments.

## Structure

### commands/
User-facing operations. Each folder represents a command that can be executed.

- start → shows execution plan
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