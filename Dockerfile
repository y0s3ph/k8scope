FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X github.com/y0s3ph/k8scope/internal/cli.Version=${VERSION} -X github.com/y0s3ph/k8scope/internal/cli.Commit=${COMMIT} -X github.com/y0s3ph/k8scope/internal/cli.Date=${DATE}" \
    -o /bin/k8scope ./cmd/k8scope

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/k8scope /usr/local/bin/k8scope

ENTRYPOINT ["k8scope"]
