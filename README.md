# Gather weather (Software Engineering School 5.0 test assessment)

## Requirements

- Go 1.24.1+
- Docker
- GNU utils (make)

## Local development

Make sure you're on Go version 1.24.1+.

Rename `.example.env` to `.env` and change environment variables you want.

### Run the project:

You can run both app and db in Docker:

```bash
docker-compose up -d
```

Or run db in Docker and app locally:

```bash
docker-compose up -d db
```

```bash
make run
```
