# frinkconv-api

`frinkconv-api` provides a JSON HTTP API for converting a value from one unit of measurement to another.

Under the hood, it uses [frinkconv](https://github.com/seanbreckenridge/frinkconv) which is a nice CLI tool
around [Frink](https://frinklang.org/) for unit of measurement conversions.

It works by spawning a configurable number of long-running `frinkconv` REPL processes at startup time, incoming HTTP requests will then pick
an available REPL from a pool (or block if none are available) with which to execute the conversion.

In my testing, it seems to be computationally fairly lightweight but there is a memory cost to running more REPL processes.

Per-request performance is not too bad; in all cases using my 2019 MacBook Pro (x86) to test:

- Docker
    - `1000 batch requests of 10 conversions each in 11.466923209s; 87.20735124615936 requests per second, 872.0735124615936 conversions per second`
- Native
    - `1000 batch requests of 10 conversions each in 4.12123353s; 242.6457983321319 requests per second, 2426.4579833213193 conversions per second`

As you'd expect, Docker for Mac is slower due to the virtualization; I daresay there's little to no cost for a Linux-based machine.

## TODO

- Functionality
    - Add support for "Multiple Conversions" ([reference](https://frinklang.org/#Conversions)); basically the process of extracting the
      destination value as a tuple
- Bugs / tech debt
    - Find a way to pin the version of `frinkconv` that's downloaded at Docker build time
    - Validate the given `source_units` and `destination_units` before attempting the conversion
        - _Possibly_ some not-yet-fully-understood risk of code injection without this
    - Probably lurking bugs in some of the REPL output parsing I'm doing
    - Look at a closer / direct integration with Frink (skip the REPL)
        - Speculated that this may be faster / less resource heavy

## Usage

### Run using published Docker image

**Prerequisites**

- [Docker](https://www.docker.com/)

**Run the server**

```shell
docker run --rm -it -p 8080:8080 initialed85/frinkconv-api
```

### Build and run using Docker Compose

**Prerequisites**

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

**Build and run the server**

```shell
docker compose up --build
```

### Build and run natively

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

## Ghetto performance testing

**Assuming you already have a `frinkconv_api` server running on port 8080:

```shell
go run test/main.go
2022/08/11 10:47:16 1000 batch requests of 10 conversions each in 3.970077757s; 251.88423532431082 requests per second, 2518.8423532431084 conversions per second
```
