# fleet-manager

A vehicle dispatch and routing service for fleet management. Built with Go, chi, and SQLite.

## Overview

fleet-manager exposes a REST API to manage vehicles, compute shortest routes between cities, and dispatch vehicles to jobs via a priority queue. Dispatch jobs are processed by a background worker that runs alongside the HTTP server.

## Architecture

```
cmd/server/         entry point, wires dependencies, starts HTTP server + dispatch worker
internal/
  vehicle/          vehicle CRUD — model, service, handler, store interface
  route/            Dijkstra shortest-path on a weighted directed graph
  dispatch/         priority queue, job model, background worker, HTTP handler
  storage/          SQLite connection, migrations, vehicle store implementation
```

## API

### Vehicles

| Method | Path | Description |
|--------|------|-------------|
| POST | `/vehicles` | Create a vehicle |
| GET | `/vehicles` | List all vehicles |
| GET | `/vehicles/{id}` | Get a vehicle |
| PATCH | `/vehicles/{id}/status` | Update vehicle status |

### Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/routes?from=X&to=Y` | Shortest path between two nodes |
| GET | `/routes/nodes` | List all known nodes in the graph |

### Dispatch

| Method | Path | Description |
|--------|------|-------------|
| POST | `/dispatch/jobs` | Enqueue a dispatch job |
| GET | `/dispatch/jobs` | Queue depth and next pending job |

## Dispatch Flow

1. Client posts a job with `from`, `to`, `urgency`, and `min_cap`.
2. Job is pushed onto a thread-safe max-heap priority queue.
3. Background worker pops the highest-priority job, picks an available vehicle meeting capacity, computes the route, and marks the vehicle `en_route`.
4. Jobs of equal urgency are processed FIFO by `created_at`.

Urgency levels: `1` low, `2` medium, `3` high, `4` critical.

## Running

```bash
# install dependencies
go mod tidy

# run the server
go run ./cmd/server/main.go
```

Server starts on `:8080`. SQLite database is created at `fleet.db` on first run.

Graceful shutdown on `SIGINT` / `SIGTERM` — the dispatch worker drains cleanly before exit.

## Example

```bash
# create a vehicle
curl -X POST http://localhost:8080/vehicles \
  -H "Content-Type: application/json" \
  -d '{"name": "Truck Alpha", "status": "available", "capacity": 10}'

# find a route
curl "http://localhost:8080/routes?from=KTM&to=JMP"

# enqueue a dispatch job
curl -X POST http://localhost:8080/dispatch/jobs \
  -H "Content-Type: application/json" \
  -d '{"from": "KTM", "to": "JMP", "urgency": 4, "min_cap": 5}'
```

## Graph

Pre-seeded with Nepali cities. Costs represent distance in km.

```
KTM -- 200km -- PKR
KTM -- 150km -- BRT
PKR -- 100km -- BRT
PKR -- 180km -- JMP
BRT -- 220km -- JMP
```

## Stack

- [Go](https://golang.org) 1.22+
- [chi](https://github.com/go-chi/chi) — HTTP router
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) — pure-Go SQLite driver
- [google/uuid](https://github.com/google/uuid) — job ID generation
