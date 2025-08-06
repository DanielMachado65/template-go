# Gin + MongoDB + Docker starter

## Run

```bash
# Build and run the API in development mode (hot reload)
docker compose up --build api-dev

# Build and run the API in production mode
docker compose up --build api
```

API will listen on http://localhost:8080

### Health

```bash
curl http://localhost:8080/healthz
```

### Create item

```bash
curl -X POST http://localhost:8080/items       -H "Content-Type: application/json"       -d '{"name":"first"}'
```

### List items

```bash
curl http://localhost:8080/items
```
