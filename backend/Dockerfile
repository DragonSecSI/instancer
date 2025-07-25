FROM golang:1.24-alpine3.22 AS builder
WORKDIR /app

RUN apk add --no-cache make

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN make dist


FROM alpine:3.22
WORKDIR /app

RUN apk add curl && \
    curl -sSf https://atlasgo.sh | sh

COPY --from=builder /app/dist/ /app/

ENTRYPOINT ["/app/entrypoint.sh"]
