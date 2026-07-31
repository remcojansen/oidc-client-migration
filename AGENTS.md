# AGENTS.md

Instructions for AI coding agents (and human contributors) working in this repository.

## Repository layout

Keep the top-level layout intentionally flat and purpose-named:

```
export/                 # Go tool: exports client configs from an auth server (Keycloak, PingFederate)
import/                 # Terraform: imports client configs into a target auth server
client-configurations/  # Hand-off artifact between export/ and import/ (canonical YAML)
docs/                   # Canonical model spec + usage guides
```

**Never introduce a new top-level directory without asking first** — even if it seems like an
obvious fit. Propose it and get confirmation before creating it.

## Commit conventions

- Use **Conventional Commits** style, with an **all-lowercase header**.
  Example: `feat: add export retry logic`, `fix: correct pingfederate secret mapping`.
- Common types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `build`.
- Do **not** add a `Co-authored-by:` trailer for the AI agent to commits in this repository.

## Branch & PR conventions

- If a branch's purpose is known in advance and clearly categorized, prefix it accordingly:
  `feat/`, `fix/`, `docs/`, `refactor/`, `test/`, `chore/`.
- If the purpose doesn't cleanly fit a category (e.g. exploratory work, mixed-purpose branch),
  use a short descriptive branch name with **no prefix** rather than forcing a category.
- Merge strategy: **squash merge** into `main` for a clean, linear history.

## Validation before considering a change "done"

Only run the validation relevant to what was touched — don't run the full matrix for unrelated changes:

- **Go changes** → `go build ./...`, `go vet ./...`, `go test ./...`
- **Terraform changes** → `terraform validate` (see `import/`)

## Dependency policy

**Always ask before adding a new dependency** — new Go modules or Terraform providers — even if
it seems clearly needed for the task. Explain why it's needed and let the user confirm first.

## Documentation upkeep

Update the relevant README/docs in the same change **only when the change affects user-facing
behavior** (new flags, new commands, changed workflow, new config options). Internal refactors
that don't change behavior don't require doc updates.

## Terraform / infrastructure safety

This project provisions real authorization server clients via Terraform. The agent must:

- **Never run `terraform apply` or `terraform destroy`** against real infrastructure.
- Only run `terraform plan` / `terraform validate` unless the user **explicitly** asks for
  `apply`/`destroy` in that specific instance.

## Secrets & credentials handling

- Never commit `.tfvars` files, tokens, credentials, or Terraform state files.
- Treat any credential-like value (client secrets, access tokens, API keys) as sensitive:
  don't print them in full in logs, commit messages, or chat output — redact or truncate them.

## Git history safety

- Never force-push (`git push --force`) or rewrite git history (`rebase -i`, amending already-
  pushed commits, etc.) unless the user **explicitly** asks for it in that specific instance.

## Temporary files

- Do not read/write directly under `/tmp`. Use a dedicated subdirectory, e.g. `/tmp/ocm-<purpose>/`,
  and clean it up when the task is done.
