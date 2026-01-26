# Product Requirements Document - build a MCP server for my daily work

## Project Overview
Build a go lang project template for microservices that I can reutilize.

---

## Phase 1: Project Setup & Core Infrastructure

- [x] Initialize a go project with the packages as a git submodule in this repository:
- [x] load configs and options using https://github.com/spf13/viper https://github.com/spf13/cobra
- [x] Write readiness probes and liveness probes in `/healthz/liveness` and `/healthz/readiness`  - readiness should fail if backend dependencies are not met
- [x] SQL backend will be sql
- [x] Write API documentation with Swagger/OpenAPI
- [x] Add a docker container with minimal dependencies and according security constrains
- [x] Configure the program to be loaded from environment variables and cli parameters
- [x] Enforce semver and conventional commits through `pre-commit` workflows
- [x] use https://github.com/golang-migrate/migrate for migrations
- [x] Set up logging framework and default to stdout and stderr
- [x] Add tests to test parsing of arguments, server working,
- [x] create the option to provide a postgresql backend
- [x] Add a docker container with minimal dependencies and according security constrains
- [x] Add tests for all the functionality implemented at this point
- [x] Create README with setup instructions
- [x] Document environment variables and configuration
- [x] Write deployment guide for production
- [x] Create user guide with screenshots
- [x] Document database schema and relationships
- [x] Add inline code comments for complex logic
- [x] Create architecture diagram

## Success Criteria

- All endpoints return proper HTTP status codes
- All tests pass with >80% coverage
- Zero critical security vulnerabilities

## Technical Stack

**Backend:**
- Go lang
- PostgreSQL
- Redis (optional)
- Kubernetes for production
- Docker-compose for local development

**DevOps:**
- Docker + Docker Compose

---

## Notes

- Use semantic versioning for releases
- Follow conventional commits for git messages
- Database migrations must be reversible
- Maintain changelog for all releases
