# Why This Repository Structure Exists: An Enterprise Argument

## Introduction

Every software organisation eventually faces the same question: how should we organise our source code? The answer shapes everything that follows — how teams collaborate, how fast features ship, how reliably systems deploy, and how painlessly new engineers onboard. Get it wrong, and you spend years fighting your own tooling. Get it right, and the structure becomes invisible — it simply works.

This document is a comprehensive argument for why this workspace is structured the way it is. It is not a quick reference card. It is a deliberate, reasoned case for every category of repository, every separation of concern, and every boundary drawn between them. If you are wondering "why not just put everything in one repo?" or "why do we need a separate deployment repo?" — this essay answers those questions in depth.

The structure described here is not invented from scratch. It draws from years of industry experience, the practices of organisations like Google, Netflix, Spotify, and Shopify, the principles of GitOps as defined by Weaveworks and the CNCF, and the hard-won lessons of teams that tried simpler approaches and discovered why they fail at scale.

---

## The Fundamental Problem: Software Is Not One Thing

A modern software system is not a single program. It is a constellation of concerns:

- **Application logic** — the code that implements business features.
- **Service topology** — how those applications connect, communicate, and depend on each other.
- **Environment configuration** — the differences between development, staging, and production.
- **Infrastructure** — the cloud resources (networks, databases, compute) that everything runs on.
- **Observability** — the dashboards, alerts, and log pipelines that tell you whether things are working.
- **Developer tooling** — the scripts, CLIs, and automation that make engineers productive.
- **Documentation** — the architecture decisions, API contracts, and operational procedures that preserve institutional knowledge.

Each of these concerns has a different lifecycle. Application code changes multiple times per day. Infrastructure changes weekly or monthly. Documentation changes when decisions are made. Alert rules change when incidents reveal gaps. If you force all of these into a single repository — or even into poorly chosen groupings — you create coupling where none should exist. A change to a Grafana dashboard triggers a CI pipeline for your payment service. A Terraform module update forces a rebuild of your API gateway. An intern updating a runbook accidentally blocks a production deployment.

The structure of your repositories is, fundamentally, a statement about which things are allowed to change independently. Every boundary you draw says: "these two concerns have different owners, different cadences, and different risk profiles, and they should not be forced to move in lockstep."

---

## Why Not a Monorepo?

The monorepo approach — putting all code in a single repository — has legitimate advocates. Google famously uses a monorepo. So do Facebook and Twitter. The argument is compelling on the surface: one repo means one source of truth, atomic cross-cutting changes, and no dependency versioning headaches.

But the monorepo argument has a critical asterisk: **it only works with custom tooling that most organisations do not have and cannot afford to build.** Google's monorepo is supported by Blaze/Bazel, Piper, CitC (Clients in the Cloud), and a dedicated infrastructure team whose sole job is making the monorepo work. Without those tools, a monorepo at scale becomes a nightmare:

- **CI/CD becomes slow.** Every commit triggers analysis of the entire repository to determine what changed. Without Google-grade build systems, this means either running everything (slow and wasteful) or maintaining fragile path-based filters (error-prone and constantly breaking).
- **Permissions become coarse.** Git does not support per-directory access control. In a monorepo, everyone can see and modify everything, or you layer on complex tooling to simulate fine-grained permissions.
- **Git performance degrades.** Git was not designed for repositories with millions of files. Clone times, checkout times, and blame operations slow to a crawl. Solutions like sparse checkout and partial clone help, but they add complexity.
- **Team autonomy suffers.** When 50 engineers share a single repository, merge conflicts become a daily occurrence. Branch protection rules apply globally. A broken CI pipeline blocks everyone, not just the team that broke it.
- **Cognitive load increases.** Opening a repository with 500 directories and 10,000 files is overwhelming. Engineers waste time navigating irrelevant code to find what they need.

The multi-repo approach trades atomic cross-cutting changes (which are rare) for independent velocity (which is constant). In an enterprise with multiple teams, multiple services, and multiple deployment targets, the multi-repo approach wins — provided the repositories are well-organised and managed by tooling that keeps them in sync.

That tooling is exactly what this workspace controller provides. The `system-definition.yaml` file is the single source of truth for which repositories exist and where they live. The controller handles cloning, pulling, and status checking across all of them. You get the discoverability of a monorepo with the independence of multi-repo.

---

## Common Repository Categories, and Why Each One Exists

The categories below are not hardcoded into the workspace controller — they are defined entirely in `system-definition.yaml`. You can add, remove, or rename categories to fit your organisation. The controller simply iterates over whatever you define. That said, these are the categories that most enterprise projects converge on, and each one exists for a specific reason.

### Applications: One Repo Per Service

The most important boundary in the entire structure is the one between services. Each microservice or application gets its own repository. This is not a stylistic preference — it is a direct consequence of the microservices architecture pattern, and violating it undermines every benefit that microservices provide.

**Independent deployment.** The entire point of microservices is that you can deploy one service without deploying all of them. If two services share a repository, their CI/CD pipelines are entangled. A commit to Service A triggers a build of Service B. A failed test in Service B blocks a deployment of Service A. The blast radius of every change expands to include code that was not changed.

**Independent versioning.** Each service has its own version number, its own release cadence, and its own changelog. When services share a repo, versioning becomes ambiguous. Does version 2.3.1 refer to the API gateway, the user service, or both? Semantic versioning loses its meaning when multiple independently-evolving components share a version.

**Team ownership.** In an enterprise, different teams own different services. Repository boundaries are the most natural and enforceable ownership boundaries. GitHub's CODEOWNERS file, branch protection rules, and access controls all operate at the repository level. When a team owns a repository, they own its CI pipeline, its deployment process, and its on-call rotation. When two teams share a repository, ownership becomes a negotiation.

**Build isolation.** Each application repo contains a Dockerfile that builds the service into a container image. The build context is the repository itself — nothing more. This guarantees that the service can be built in isolation, without access to other services' source code. If a service cannot build without another service's code, it is not a microservice — it is a module in a monolith pretending to be distributed.

**The rule is simple:** an application repo builds and publishes a container image. It never references other service repos directly. If two services need to communicate, they do so through APIs, message queues, or shared contracts — never through shared source code.

---

### Orchestrator: Wiring Services Together

If each service is an independent unit, something needs to define how they connect. That is the orchestrator's job.

The orchestrator repository contains Docker Compose files (or Kubernetes manifests, or Helm charts) that pull pre-built images and wire them together with networking, ports, volumes, and dependency ordering. It is the answer to the question every developer asks on their first day: "How do I run the whole thing locally?"

**Why a separate repo?** Because the service topology is not owned by any single service. It is a cross-cutting concern. If you put the docker-compose file in one of the application repos, you create an asymmetry — one service becomes "special" because it contains the orchestration config. Changes to the topology require commits to that service's repo, polluting its history with unrelated changes.

**Why not build from source?** The orchestrator references pre-built images, not source code. This is a critical distinction. If the orchestrator builds services from source, it needs access to every service's repository, and a change to any service's code triggers a rebuild of the entire stack. By referencing images, the orchestrator is decoupled from the build process. It only cares about which version of each service to run, not how they are built.

This separation also mirrors production. In production, you never build from source — you deploy pre-built, tested, versioned images. The orchestrator should work the same way in development, so that the local environment is as close to production as possible. The only difference is that developers might use `:latest` or `:dev` tags locally, while production uses pinned versions.

---

### Deployment: Environment Configuration

Here is a question that trips up many teams: where do environment-specific settings live? The database connection string for staging is different from production. The feature flags for development are different from staging. The replica count, the resource limits, the TLS certificates — all of these vary by environment.

The deployment repository is the answer. It contains a folder for each environment (`dev/`, `staging/`, `prod/`), and each folder contains the configuration values specific to that environment.

**Why not branches?** The branch-per-environment pattern — where `main` is production, `staging` is staging, and `develop` is development — is one of the most persistent anti-patterns in the industry. It fails for three reasons:

1. **Drift.** Over time, the branches diverge. A hotfix applied to `main` is not cherry-picked to `staging`. A config change in `develop` is forgotten when promoting to `staging`. The environments become inconsistent, and no one knows which branch reflects reality.

2. **Merge conflicts.** Environment-specific values conflict when merging between branches. The staging database URL conflicts with the production database URL. These conflicts are not real conflicts — they are artefacts of using branches for something branches were not designed for.

3. **Auditability.** With folders, you can see all environments in a single commit. You can diff `dev/values.yaml` against `prod/values.yaml` and immediately see the differences. With branches, comparing environments requires checking out different branches and manually diffing — a process that is error-prone and rarely done.

The folder-per-environment approach is recommended by the ArgoCD project, the Flux project, Codefresh, and virtually every GitOps practitioner who has tried both approaches. It is the industry consensus.

**Why separate from the orchestrator?** Because the orchestrator defines *what* runs (which services, how they connect), while the deployment repo defines *how* it runs in each environment (how many replicas, which database, what resource limits). These are different concerns with different change cadences. The orchestrator changes when you add or remove a service. The deployment config changes when you tune performance, rotate secrets, or promote a release.

---

### Infrastructure: What Exists in the Cloud

Infrastructure-as-Code (IaC) is no longer optional. If your cloud resources are not defined in version-controlled code, you are operating on hope and tribal knowledge. The infrastructure repository contains Terraform modules (or Pulumi, CloudFormation, or CDK) that define the cloud resources your system needs: VPCs, subnets, databases, storage buckets, Kubernetes clusters, DNS records, and everything else.

**Why separate from application repos?** Because infrastructure has a fundamentally different lifecycle. Application code changes daily. Infrastructure changes weekly or monthly. The people who modify Terraform modules (platform engineers, SREs) are often different from the people who write application code (product engineers). The review process is different — infrastructure changes require careful planning, blast radius analysis, and often a change management process.

If infrastructure code lives in an application repo, it sends the wrong signal. It suggests that infrastructure changes are as routine as application changes. It makes it easy for application developers to accidentally modify infrastructure. It clutters the application repo's history with Terraform state changes and module updates that have nothing to do with the application.

**Why separate from deployment config?** Because infrastructure defines *what exists* (a database cluster, a Kubernetes namespace, a load balancer), while deployment config defines *what runs on it* (which services, with what settings). You provision infrastructure once (or rarely). You deploy applications constantly. These cadences should not be coupled.

The infrastructure repo defines the stage. The deployment repo directs the actors. The application repos are the actors themselves.

---

### Observability: Knowing What Is Happening

You cannot operate what you cannot observe. The observability repository contains everything related to monitoring, logging, and alerting: Grafana dashboard definitions, Prometheus recording and alerting rules, log pipeline configurations, and on-call escalation policies.

**Why a separate repo?** Because observability is a cross-cutting concern that spans all services. A single Grafana dashboard might display metrics from five different services. An alerting rule might fire based on the error rate of the API gateway combined with the latency of the database. If these definitions live in individual service repos, there is no single place to see the complete observability picture.

Additionally, observability configuration changes for different reasons than application code. You add a dashboard after an incident reveals a blind spot. You tune an alert threshold after too many false positives. You add a log pipeline when a new compliance requirement emerges. None of these changes are related to application features, and they should not pollute application repos' histories.

Centralising observability also enables consistency. When all dashboards are in one repo, you can enforce naming conventions, ensure every service has baseline metrics, and review alerting rules holistically. When dashboards are scattered across 20 service repos, inconsistency is inevitable.

---

### Platform: The Developer Experience Layer

The platform repository contains internal developer tooling — the tools that make engineers productive but are not themselves part of the product. This workspace controller is the primary example: it is a CLI that manages all the other repositories, but it is not a microservice that serves end users.

**Why separate from tools?** Because platform tooling is foundational. It is the layer that everything else depends on. The workspace controller defines how repositories are cloned, updated, and validated. CI/CD template libraries define how every service is built and deployed. These are not utilities — they are infrastructure for the development process itself.

Platform repos tend to have different ownership (a platform or DevEx team), different SLAs (if the workspace controller breaks, everyone is blocked), and different release processes (changes need to be backwards-compatible because every team depends on them).

---

### Tools: Everything Else That Helps

Every project accumulates small utilities that do not fit neatly into any other category. Migration scripts. Data seeders. Code generators. One-off automation for bulk operations. If these do not have a home, they end up in random places — a `scripts/` folder in an application repo, a personal repository on someone's GitHub, a Slack message that gets lost.

The tools repository is the designated home for these utilities. It is intentionally broad because its purpose is to prevent clutter elsewhere. The rule is simple: if a utility serves multiple repos or does not belong to a specific service, it goes in tools.

This category also includes deployment helpers and automation scripts. CI/CD pipeline definitions might reference shared scripts from the tools repo. A deployer utility that wraps Terraform commands might live here. The key distinction from platform is scope: platform tools are foundational and used by everyone; tools are useful but not critical.

---

### Docs: Institutional Knowledge

Code is not documentation. Comments explain *how*; documentation explains *why*. Architecture Decision Records capture the reasoning behind choices that will otherwise be forgotten. API contracts define the agreements between services. Manual runbooks describe the procedures that humans follow during incidents.

**Why a separate repo?** Because documentation that spans multiple services has no natural home in any single service repo. An ADR about choosing PostgreSQL over MongoDB affects every service that uses the database. An API contract between the user service and the notification service belongs to neither. An onboarding guide that walks a new engineer through the entire system cannot live in one service's README.

**Why only manual runbooks?** Automated runbooks — scripts that remediate incidents automatically — belong in the service or infrastructure repo where they execute. A script that restarts a crashed service belongs in that service's repo. A Terraform module that scales up a database belongs in the infrastructure repo. The docs repo is for procedures that require human judgement: "If the payment service is returning 500 errors, check the following three things before escalating to the payments team."

---

## The Deeper Principles

### Single Responsibility, Applied to Repositories

The Single Responsibility Principle states that a module should have one, and only one, reason to change. This principle applies to repositories just as it applies to classes and functions.

An application repo changes when business logic changes. The orchestrator changes when service topology changes. The deployment repo changes when environment configuration changes. The infrastructure repo changes when cloud resources change. Each repository has one reason to change, one type of owner, and one cadence of evolution.

When a repository has multiple reasons to change, it becomes a coordination bottleneck. Two teams need to merge to the same repo. Two CI pipelines compete for the same branch. Two types of changes — one urgent, one routine — are forced through the same review process. Single responsibility at the repository level prevents this.

### Conway's Law, Embraced Rather Than Fought

Conway's Law states that organisations design systems that mirror their communication structures. This is often cited as a problem, but it is better understood as a force of nature. You cannot fight it — you can only align with it.

The repository structure mirrors the team structure of an enterprise:

- **Product teams** own application repos.
- **A platform team** owns the platform and tools repos.
- **An SRE or infrastructure team** owns the infrastructure and observability repos.
- **A DevOps or release team** owns the orchestrator and deployment repos.
- **Everyone** contributes to docs.

When repository boundaries align with team boundaries, ownership is clear, communication overhead is minimised, and each team can move at its own pace. When they do not align — when two teams share a repo, or one team's code is scattered across five repos — friction is constant.

### The Principle of Least Coupling

Coupling is the degree to which one component depends on another. In a well-designed system, coupling is minimised — components interact through well-defined interfaces and can be changed independently.

The repository structure enforces low coupling at the highest level. Application repos do not reference each other. The orchestrator references images, not source code. The deployment repo references configuration values, not application logic. The infrastructure repo provisions resources without knowing which applications will use them.

This is not accidental. Every boundary in the structure is drawn at a point where coupling should be minimal. If you find yourself needing to change two repos simultaneously for a single feature, it is a signal that either the feature crosses a legitimate boundary (and the coordination cost is acceptable) or the boundary is drawn in the wrong place (and the structure should be reconsidered).

### GitOps: The Repository as the Source of Truth

GitOps is the practice of using Git as the single source of truth for both application code and operational configuration. The repository structure is designed with GitOps in mind:

- The **deployment repo** is the source of truth for what should be running in each environment. A GitOps controller (ArgoCD, Flux) watches this repo and reconciles the actual state of the cluster with the desired state defined in Git.
- The **infrastructure repo** is the source of truth for what cloud resources should exist. A CI pipeline applies Terraform changes when the repo is updated.
- The **observability repo** is the source of truth for what should be monitored. Dashboards and alert rules are provisioned from Git, not created manually in a UI.

When every operational concern is defined in a Git repository, you get auditability (who changed what, when, and why), rollback capability (revert a commit to undo a change), and review processes (pull requests for infrastructure changes). These properties are not possible when configuration lives in UIs, wikis, or engineers' heads.

### Blast Radius Minimisation

Every change carries risk. The goal of good architecture is to minimise the blast radius of any single change — to ensure that a mistake in one area does not cascade into failures in unrelated areas.

Repository boundaries are blast radius boundaries. A broken CI pipeline in one application repo does not block other application repos. A misconfigured Terraform module does not prevent application deployments. A typo in a Grafana dashboard does not trigger a rebuild of the payment service.

This is why the structure has many focused categories rather than two or three. Each additional boundary reduces the blast radius of changes within that category. The cost is coordination overhead when changes cross boundaries — but cross-boundary changes are rare by design, and the coordination cost is a feature, not a bug. It forces you to think carefully about changes that affect multiple concerns.

---

## The Workspace Controller: Making It Work

A multi-repo structure without tooling is a multi-repo mess. Engineers forget which repos exist. New team members do not know what to clone. Repositories drift out of sync. The workspace controller solves this by providing a single command to clone all repositories, check their status, and keep them updated.

The `system-definition.yaml` file is the manifest. It lists every repository, organised by category, with its Git URL and local path. The controller reads this manifest and performs operations across all repositories:

- `init` — scaffolds a new workspace with a definition file, Makefile, and .gitignore.
- `clone` — clones all repositories defined in the manifest.
- `pull` — pulls the latest changes for all repositories (clones missing ones).
- `push` — pushes local commits for all repositories.
- `checkout` — switches all repositories to their defined branch (or an override).
- `status` — shows the branch, dirty/clean state, and ahead/behind counts for every repository.
- `validate` — verifies that every repository matches the definition.
- `doctor` — diagnoses the local environment (Git, Go, SSH, Docker).
- `ssh` — interactive wizard for SSH key setup.

This is the key insight that makes multi-repo viable: **you do not need a monorepo if you have a tool that treats multiple repos as a coherent workspace.** The workspace controller provides the discoverability and convenience of a monorepo while preserving the independence and isolation of multi-repo.

The manifest is also documentation. A new engineer can read `system-definition.yaml` and immediately understand the full scope of the system — every service, every infrastructure component, every tool. This is something that a monorepo provides implicitly (everything is in one place) but that multi-repo usually lacks (you have to know which repos exist). The manifest bridges this gap.

---

## Common Objections, Addressed

### "This is too many repos."

The number of categories does not equal the number of repos. The applications category alone might contain 20 repos (one per service). The total count depends on the size of the system. The categories themselves are dynamic — defined in `system-definition.yaml`, not hardcoded in the tooling. You can add, remove, or rename them freely. But the common categories described here are minimal — removing any one of them forces its contents into another category where it does not belong, creating the coupling and confusion that the structure is designed to prevent.

### "Cross-cutting changes are harder."

Yes. If you need to change an API contract, update the service that implements it, modify the orchestrator to add a new dependency, and update the deployment config for the new service — that is four repos. In a monorepo, it would be one commit.

But this friction is intentional. Cross-cutting changes are high-risk changes. They affect multiple teams, multiple deployment pipelines, and multiple environments. Forcing them through separate pull requests in separate repos ensures that each change is reviewed by the appropriate team, tested in isolation, and deployed independently. The alternative — a single commit that changes everything at once — is convenient but dangerous.

### "We are a small team. We do not need this."

You might not need all of these categories today. But the structure scales gracefully. Start with applications, orchestrator, and infrastructure. Add deployment when you have multiple environments. Add observability when you have dashboards to manage. Add docs when you have architecture decisions to record. The categories are independent — you can adopt them incrementally.

What you should not do is start with a monorepo and try to split it later. Repository splits are painful, disruptive, and often incomplete. Starting with the right structure — even if some categories are empty — is far cheaper than restructuring later.

### "Our CI/CD is simpler with a monorepo."

Simpler, perhaps, but also more fragile. A monorepo CI pipeline must determine which parts of the codebase changed and run only the relevant tests and builds. This requires path-based filtering, which is brittle and error-prone. When it breaks, everything breaks.

With multi-repo, each repository has its own CI pipeline. The pipeline is simple because it only builds one thing. There is no path filtering, no conditional logic, no "skip if only docs changed" heuristics. The pipeline runs on every commit and builds the one artefact that the repo produces. This simplicity is a feature.

---

## Conclusion

The repository structure described in this document is not arbitrary. Every category exists because the concern it represents has a different lifecycle, a different owner, and a different risk profile from every other category. Every boundary is drawn at a point where coupling should be minimal and independence should be maximal.

The structure follows established principles: single responsibility, Conway's Law, least coupling, blast radius minimisation, and GitOps. It aligns with the practices of mature engineering organisations and the recommendations of the GitOps community.

It is not the only valid structure. But it is a structure that has been tested by industry, refined through experience, and designed to scale from a small team to a large enterprise. The workspace controller makes it practical by providing the tooling that multi-repo requires to be manageable.

The question is not "why so many repos?" The question is "why would you force things that change independently into the same repo?" The answer to the second question is almost always "because we did not have the tooling to manage multiple repos." This workspace provides that tooling. The structure follows naturally.

---

## Appendix: Repository Category Quick Reference

| Category | Purpose | Changes When | Typical Owner |
|---|---|---|---|
| Applications | Business logic, one repo per service | Features are added or bugs are fixed | Product teams |
| Orchestrator | Wires services together for local/full-stack runs | Services are added, removed, or re-wired | Platform / DevOps team |
| Deployment | Per-environment config (dev/staging/prod folders) | Environment settings change, releases are promoted | Release / DevOps team |
| Infrastructure | Cloud resource provisioning (Terraform, IaC) | Cloud resources are added or modified | SRE / Infrastructure team |
| Observability | Monitoring, dashboards, alerts, log pipelines | Incidents reveal blind spots, thresholds are tuned | SRE / Observability team |
| Platform | Internal developer tooling (e.g. workspace controller) | Developer workflows change | Platform / DevEx team |
| Tools | General-purpose utilities and scripts | New automation needs arise | Any team |
| Docs | Architecture docs, API contracts, manual runbooks | Decisions are made, procedures change | Any team |
