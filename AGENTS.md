## MANDATORY: No CI/CD Pipelines

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
