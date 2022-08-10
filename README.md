# frinkconv-api

`frinkconv-api` provides a JSON HTTP API for converting a value from one unit of measurement to another.

Under the hood, it uses [frinkconv](https://github.com/seanbreckenridge/frinkconv) which is a nice CLI tool
around [Frink](https://frinklang.org/) for unit of measurement conversions.

It works by spawning a configurable number of long-running `frinkconv` REPL processes at startup time, incoming HTTP requests will then pick
an available REPL from a pool (or block if none are available) with which to execute the conversion.

In my testing, it seems to be computationally fairly lightweight but there is a memory cost to running more REPL processes.

Individual request times aren't fantastic, so the sweet spot for high volume systems will be somewhere between batching the requests on the
client side and running an appropriate amount of REPL processes.

## TODO

- There may be an opportunity for a closer integration with Frink (which might be more responsive and not as heavy)
- Definitely there are lurking bugs in some of the fragile REPL output parsing I'm doing

## Usage

### Native

**Prerequisites**

- [Go 1.17](https://go.dev/)

**Build the server**

```shell
go build -o bin/frinkconv-api cmd/main.go
```

**Run the server**

```shell
bin/frinkconv-api -port 8080 -processes 4
```

### Docker Compose

**Prerequisites**

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

**Build the server**

```shell
docker compose build
```

**Run the server**

```shell
docker compose up
```

## Interactions via the API

**Prerequisites**

- [Curl](https://github.com/curl/curl)
- [jq](https://github.com/stedolan/jq)

**Attempt some single conversions**

```shell
# this should succeed
curl -s -X POST -d '{"source_value": 120.0, "source_units": "feet", "destination_units": "metres"}' http://localhost:8080/convert/ | jq
{
  "destination_value": 36.576
}

# this should fail
curl -s -X POST -d '{"source_value": 120.0, "source_units": "apples", "destination_units": "oranges"}' http://localhost:8080/convert/ | jq
{
  "error": "Warning: undefined symbol \"apples\".\nUnknown symbol \"oranges\"\nWarning: undefined symbol \"apples\".\nWarning: undefined symbol \"oranges\".\nUnconvertable expression:\n  120 apples (undefined symbol) -> oranges (undefined symbol)"
}
```

**Attempt a batch conversions**

```shell
# this should succeed and fail
curl -s -X POST -d '[{"source_value": 120.0, "source_units": "feet", "destination_units": "metres"}, {"source_value": 120.0, "source_units": "apples", "destination_units": "oranges"}, {"source_value": 120.0, "source_units": "feet", "destination_units": "metres"}, {"source_value": 120.0, "source_units": "apples", "destination_units": "oranges"}]' http://localhost:8080/batch_convert/ | jq
[
  {
    "destination_value": 36.576
  },
  {
    "error": "Warning: undefined symbol \"apples\".\nUnknown symbol \"oranges\"\nWarning: undefined symbol \"apples\".\nWarning: undefined symbol \"oranges\".\nUnconvertable expression:\n  120 apples (undefined symbol) -> oranges (undefined symbol)"
  },
  {
    "destination_value": 36.576
  },
  {
    "error": "Warning: undefined symbol \"apples\".\nUnknown symbol \"oranges\"\nWarning: undefined symbol \"apples\".\nWarning: undefined symbol \"oranges\".\nUnconvertable expression:\n  120 apples (undefined symbol) -> oranges (undefined symbol)"
  }
]
```
