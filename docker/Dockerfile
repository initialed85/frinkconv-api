FROM golang:1.17 AS builder

WORKDIR /srv/

COPY go.mod /srv/go.mod
COPY go.sum /srv/go.sum

# TODO: not supported by GitHub Actions
# RUN --mount=type=cache,target=/root/.cache/go-build go mod download

RUN go mod download

COPY cmd /srv/cmd
COPY internal /srv/internal
COPY pkg /srv/pkg

# TODO: not supported by GitHub Actions
# RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o frinkconv-api /srv/cmd/main.go

RUN go build -v -o frinkconv-api /srv/cmd/main.go

FROM golang:1.17 AS runner

RUN apt-get update && apt-get install -y git default-jre rlwrap wget

WORKDIR /srv/

RUN git clone https://github.com/seanbreckenridge/frinkconv.git

WORKDIR /srv/frinkconv

RUN ./frinkconv

RUN cp -frv ./frinkconv /usr/local/bin

WORKDIR /srv/

COPY --from=builder /srv/frinkconv-api /srv/frinkconv-api

# TODO
# HEALTHCHECK --interval=1s --timeout=1s --start-period=60s CMD ["curl", "-f", "http://localhost:8080/healthz"]

STOPSIGNAL SIGINT

# note: ensure to use JSON form (square brackets) otherwise Docker wraps CMD w/ bash (which breaks STOPSIGNALs)
ENTRYPOINT ["/srv/frinkconv-api"]
CMD ["-port", "8080", "-processes", "1"]
