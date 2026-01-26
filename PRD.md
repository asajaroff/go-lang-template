# Product Requirements Document - build a MCP server for my daily work

## Project Overview
Build a go lang project template for microservices that I can reutilize.

---

## Phase 1: Project Setup & Core Infrastructure

- [x] Initialize a go project with the packages as a git submodule in this repository:
- [x] load configs and options using https://github.com/spf13/viper https://github.com/spf13/cobra
- [x] Write readiness probes and liveness probes in `/healthz/liveness` and `/healthz/readiness`  - readiness should fail if backend dependencies are not met
- [ ] SQL backend will be sql
- [x] Write API documentation with Swagger/OpenAPI
- [ ] Add a docker container with minimal dependencies and according security constrains
- [ ] Configure the program to be loaded from environment variables and cli parameters
- [ ] Enforce semver and conventional commits through `pre-commit` workflows
- [ ] use https://github.com/golang-migrate/migrate for migrations
- [ ] Set up logging framework and default to stdout and stderr
- [ ] Add tests to test parsing of arguments, server working, 
- [ ] create the option to provide a postgresql backend
- [ ] Add a docker container with minimal dependencies and according security constrains
- [ ] Add tests for all the functionality implemented at this point
- [ ] Create README with setup instructions
- [ ] Document environment variables and configuration
- [ ] Write deployment guide for production
- [ ] Create user guide with screenshots
- [ ] Document database schema and relationships
- [ ] Add inline code comments for complex logic
- [ ] Create architecture diagram

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
