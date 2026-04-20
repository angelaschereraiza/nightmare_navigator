FROM golang:1.21.8 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/nightmare_navigator ./cmd/nightmare_navigator

FROM registry.fedoraproject.org/fedora-minimal:latest
WORKDIR /app

RUN microdnf install -y ca-certificates && microdnf clean all

COPY --from=builder /out/nightmare_navigator /usr/local/bin/nightmare_navigator

ENTRYPOINT ["/usr/local/bin/nightmare_navigator"]