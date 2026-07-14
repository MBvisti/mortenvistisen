---
name: andurel
description: Use this skill for Andurel framework projects when deciding where code belongs, adding or changing resources, controllers, models, services, routes, templ views, Inertia/Vue screens, background jobs, migrations, config, clients, or framework-adjacent internals. Enforces CLI-first generation and updates, including command discovery, dry runs, artifact inspection, and verification before manual edits.
---

# Andurel

Use this skill when working in an Andurel project or generating Andurel code. It helps place code in the right layer and use the `andurel` CLI safely.

## Mandatory CLI-First Generation

Using the `andurel` CLI is required whenever a CLI command exists for creating, updating, synchronizing, or checking an artifact. This is an invariant, not a preference. A manually correct result does not excuse skipping the CLI workflow.

This rule covers both new artifacts and changes to existing generated resources, including models, factories, controllers, scaffolds, jobs, routes, emails, and any future generator exposed by the CLI.

For every generation or generated-resource update, perform this sequence in order:

1. Run `andurel project info --json`.
2. Run `andurel commands --json`, then inspect the exact command with `andurel <command> --help`.
3. Run the exact generator or updater with `--dry-run --diff --json` when those flags are supported.
4. Inspect the structured result, especially created, updated, skipped, and command artifact arrays.
5. Run the same CLI command without `--dry-run` to apply the change. Use non-interactive confirmation flags where available.
6. Add custom behavior only after generation, and only outside CLI-managed regions unless the command explicitly supports the edit.
7. Run the relevant CLI `--check`, sync, route generation, and `andurel doctor --json` commands after customization.

Do not create or update a generated artifact by hand before completing CLI discovery. Manual creation is allowed only when command discovery proves that no applicable CLI generator or updater exists, or the available command cannot represent the required artifact. When falling back to manual work:

1. State the missing CLI capability explicitly in the work update.
2. Follow the closest generated project pattern.
3. Run any related CLI check, sync, or doctor command afterward.

If a generated file and formatter disagree, leave CLI-managed regions in the canonical form expected by the CLI and rerun the relevant `--check` command last.

## Agent Invariants

- Always use `andurel --agent --help` and `andurel commands --json` for discovery before generation.
- Run `andurel project info --json` before generation.
- Use `--json` or `--jq` when extracting data.
- Use `--dry-run --diff --json` before every generator or generated-resource update when supported.
- Inspect returned artifact arrays before assuming which files changed.
- After adding or changing Inertia routes, run `andurel generate routes --json` so frontend pages can import `resources/js/routes.ts`.
- Follow the repository rules for verification.
- Prefer the local project pattern over a generic Rails, Echo, Bun, Templ, or Vue convention.
- Keep controllers as HTTP adapters: parse input, call models or services, map errors, and render a response.
- Create a service only when there is real application orchestration, not just because code exists.

## Read When Placing Code

Read [references/layer-placement.md](references/layer-placement.md) before adding or moving behavior across models, services, controllers, routes, views, queue jobs, config, clients, or internal packages.

## First Pass

1. Inspect the existing resource closest to the requested change.
2. Identify the delivery surface: public hypermedia page, admin Inertia page, API endpoint, background job, email, or CLI/command.
3. Identify the domain object or workflow being changed.
4. Keep changes in the smallest layer that can own the behavior honestly.
5. Complete the mandatory CLI-first sequence before generating or updating generated project files.

## Layer Placement

- Put invariant business rules, domain validation, entity construction, persistence methods, and finder/query methods in `models/`.
- Put test factory definitions and factory helpers in `models/factories/`.
- Put transactions, cross-model coordination, external side effects, and multi-step application workflows in `services/`.
- Put HTTP-specific concerns in `controllers/`, `controllers/admin/`, or `controllers/api/`.
- Put route names, route paths, and URL builders in `router/routes/`.
- Put templ rendering helpers and presentation-specific adapters in `views/`.
- Put admin Inertia pages and reusable Vue components in `resources/js/`.
- Put River job argument types in `queue/jobs/` and worker implementations or registration in `queue/`.
- Put provider adapters in `clients/`, email templates/helpers in `email/`, and config/environment loading in `config/`.
- Put reusable framework-like support that is independent of one resource in `internal/`.
- Register new constructors in the existing `fx` modules for the package that owns them.

## Output Modes

Use structured output by default when automating:

| Flag | Use |
|------|-----|
| `--json` | Full `{ok,data,summary,breadcrumbs}` envelope |
| `--agent` | Structured output with non-essential human progress suppressed |
| `--jq '.field.path'` | Built-in simple field-path extraction |
| `--quiet` | Suppress human-only output |
| `--md` | Markdown output where supported |

Structured failures include `ok:false`, a stable `code`, `error`, optional `hint`, and `exit_code`. Prefer the `hint` and `breadcrumbs` fields over guessing the next command.

## Common Workflows

Inspect a project:

```bash
andurel project info --json
andurel routes --json
andurel models --json
andurel migrations --json
andurel commands --json
```

Preview scaffold generation:

```bash
andurel generate scaffold Product --dry-run --diff --json
```

Generate and review artifacts:

```bash
andurel generate scaffold Product --json
```

Update an existing model and its generated factory:

```bash
andurel generate model Newsletter --update --dry-run --diff --json
andurel generate model Newsletter --update --yes --json
andurel generate factory Newsletter --check --json
```

Generate a background job:

```bash
andurel generate job DeleteNewsletterCover --dry-run --diff --json
andurel generate job DeleteNewsletterCover --json
```

Generate Inertia route helpers:

```bash
andurel routes --json
andurel generate routes --json
```

`andurel generate routes` reads `router/routes/*.go` as the source of truth and writes `resources/js/routes.ts`. It only runs when `andurel.lock` has `scaffoldConfig.inertia` set to `vue` or `react`. Import helpers from that file in Vue or React pages instead of hard-coding URLs.

Check or sync factories:

```bash
andurel generate factory Product --check --json
andurel generate factory Product --sync --json
andurel generate factories --check --json
andurel generate factories --sync --json
```

Factory guidance:

1. Treat model `Entity` structs as the source of truth for generated factory fields.
2. Keep reusable test data builders in `models/factories/`.
3. Always run `andurel generate factory NAME --check --json` before considering a factory edit.
4. Use `--sync` to update Andurel generated regions and preserve custom helpers outside those regions.
5. Pass `--skip-factory` only when a generated model or scaffold should intentionally omit a factory.

Generate a named database seed:

1. Inspect the relevant models and existing factories in `models/factories`.
2. Add a seed function to `database/seeds`, using only exported model/factory/storage primitives.
3. Register it in `seeds.Registry` with a stable lowercase name.
4. Keep the seed idempotence expectations explicit in code comments when it may be re-run.
5. Verify the seed is discoverable:

```bash
andurel database seed --list
andurel database seed development
andurel database seed test
```

Check project health:

```bash
andurel doctor --json
```

In Inertia projects, `doctor` checks whether `resources/js/routes.ts` matches the current `router/routes/*.go` manifest. If the `routes.ts` check fails, run `andurel generate routes --json`.

## Validation

Use the repository's allowed validation commands and project guidance. In this repo, do not run `go test`, `go build`, or `npm run`; use `go vet`, `go fix`, and `gofmt`.
