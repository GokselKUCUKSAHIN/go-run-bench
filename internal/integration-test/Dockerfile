FROM golang:1.22-bookworm

WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN go build -o /usr/local/bin/go-run-bench .

COPY internal/integration-test/fixture /fixture

WORKDIR /work
