## MANDATORY: No CI/CD Pipelines

## MANDATORY: Project-Agnostic / 100% Decoupled

**This module MUST remain 100% decoupled from any consuming project. It is designed for generic use with ANY project, not one specific consumer.**

- NEVER hardcode project-specific package names, endpoints, device serials, or region-specific data
- NEVER import anything from a consuming project
- NEVER add project-specific defaults, presets, or fixtures into source code
- All project-specific data MUST be registered by the caller via public APIs — never baked into the library
- Default values MUST be empty or generic

Violations void the release. Refactor to restore generic behaviour before any commit.

**NO GitHub Actions, GitLab CI/CD, or any automated pipeline may exist in this repository!**

- No `.github/workflows/` directory
- No `.gitlab-ci.yml` file
- No Jenkinsfile, .travis.yml, .circleci, or any other CI configuration
- All builds and tests are run manually or via Makefile targets
- This rule is permanent and non-negotiable

## Agent Instructions

- Run `go build ./...` before committing
- Run `go test ./pkg/models/... ./pkg/retry/... ./pkg/circuit/... ./pkg/health/... -race -count=1` for core tests
- Provider tests that make real API calls may fail without API keys - this is expected
- Never create CI/CD configuration files


## ⚠️ MANDATORY: NO SUDO OR ROOT EXECUTION

**ALL operations MUST run at local user level ONLY.**

This is a PERMANENT and NON-NEGOTIABLE security constraint:

- **NEVER** use `sudo` in ANY command
- **NEVER** use `su` in ANY command
- **NEVER** execute operations as `root` user
- **NEVER** elevate privileges for file operations
- **ALL** infrastructure commands MUST use user-level container runtimes (rootless podman/docker)
- **ALL** file operations MUST be within user-accessible directories
- **ALL** service management MUST be done via user systemd or local process management
- **ALL** builds, tests, and deployments MUST run as the current user

### Container-Based Solutions
When a build or runtime environment requires system-level dependencies, use containers instead of elevation:

- **Use the `Containers` submodule** (`https://github.com/vasic-digital/Containers`) for containerized build and runtime environments
- **Add the `Containers` submodule as a Git dependency** and configure it for local use within the project
- **Build and run inside containers** to avoid any need for privilege escalation
- **Rootless Podman/Docker** is the preferred container runtime

### Why This Matters
- **Security**: Prevents accidental system-wide damage
- **Reproducibility**: User-level operations are portable across systems
- **Safety**: Limits blast radius of any issues
- **Best Practice**: Modern container workflows are rootless by design

### When You See SUDO
If any script or command suggests using `sudo` or `su`:
1. STOP immediately
2. Find a user-level alternative
3. Use rootless container runtimes
4. Use the `Containers` submodule for containerized builds
5. Modify commands to work within user permissions

**VIOLATION OF THIS CONSTRAINT IS STRICTLY PROHIBITED.**


