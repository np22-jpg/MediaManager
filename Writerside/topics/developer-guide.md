# Developer Guide

This section is for those who want to contribute to Media Manager or understand its internals.

### Source Code

- `media_manager/`: Backend FastAPI application (Python)
- `web/`: Frontend SvelteKit application (TypeScript)

### Backend Development

- Uses `uv` for dependency management (see `pyproject.toml` and `uv.lock`)
- Follows standard FastAPI project structure
- Database migrations are handled by Alembic (`alembic.ini`, `alembic/` directory)

### Frontend Development

- Uses `npm` for package management (see `web/package.json`)
- SvelteKit with TypeScript

### Contributing

- Please refer to the project's GitHub repository for contribution guidelines (e.g., forking, branching, pull requests)
- Consider opening an issue to discuss significant changes before starting work

## Tech Stack

### Backend

- **Framework:** Python with FastAPI
- **Database ORM:** SQLAlchemy
- **Database Migrations:** Alembic
- **Dependency Management:** uv

### Frontend

- **Framework:** SvelteKit
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **Components:** shadcn-svelte for UI components

### Deployment & CI/CD

- Docker & Docker Compose
- GitHub Actions