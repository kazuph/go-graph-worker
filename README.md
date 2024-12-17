# go-graph-worker

A Cloudflare Worker that generates line graphs from numerical data using Go.
This worker provides a simple HTTP API to create PNG graph images from JSON data arrays.

## Notice

- A free plan Cloudflare Workers only accepts ~1MB sized workers.
  - Go Wasm binaries easily exceeds this limit, so **you'll need to use a paid plan of Cloudflare Workers** (which accepts ~5MB sized workers).
- This project uses the [`workers`](https://github.com/syumai/workers) package and [gonum/plot](https://github.com/gonum/plot) for graph generation.

## Features

This worker provides the following endpoints:

- `/` (GET): A simple test endpoint that returns "Hello!"
- `/plot` (POST): Generates a line graph from the provided JSON data array
  - Accepts: JSON object with the following structure:
    ```json
    {
      "data": [number],        // Required: Array of numbers
      "title": string,         // Optional: Graph title
      "x_label": string,      // Optional: X-axis label
      "y_label": string       // Optional: Y-axis label
    }
    ```
  - Returns: PNG image of the generated graph
  - The X-axis represents the index of each data point
  - The Y-axis represents the value of each data point
  - Default values if not specified:
    - Title: "Example Plot"
    - X Label: "X"
    - Y Label: "Y"

## Requirements

- Node.js
- [wrangler](https://developers.cloudflare.com/workers/wrangler/)
  - just run `npm install -g wrangler`
- Go 1.21.0 or later

## Getting Started

* If not already installed, please install the [gonew](https://pkg.go.dev/golang.org/x/tools/cmd/gonew) command.

```console
go install golang.org/x/tools/cmd/gonew@latest
```

* Create a new project using this template.
  - Second argument passed to `gonew` is a module path of your new app.

```console
gonew github.com/syumai/workers/_templates/cloudflare/worker-go your.module/my-app # e.g. github.com/syumai/my-app
cd my-app
go mod tidy
make dev # start running dev server
```

- To change worker name, please edit `name` property in `wrangler.toml`.

## Development

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy  # deploy worker
```

### Testing dev server

Test the root endpoint:
```
$ curl http://localhost:8787/
Hello!
```

Generate a graph using the plot endpoint:

Basic usage (data only):
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"data": [1.0, 2.5, 1.8, 3.2, 4.0]}' \
  http://localhost:8787/plot >| graph.png
```

With custom title and labels:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "data": [1.0, 2.5, 1.8, 3.2, 4.0],
    "title": "Temperature Graph",
    "x_label": "Time",
    "y_label": "Temperature (°C)"
  }' \
  http://localhost:8787/plot >| custom_graph.png
```

This will create PNG files containing line graphs of the provided data points with the specified customizations.
