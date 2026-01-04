# syntax=docker/dockerfile:1

ARG PORT=8080

FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /api cmd/api/main.go
RUN go build -o /worker cmd/worker/main.go

FROM alpine:3.23 AS api
WORKDIR /
COPY --from=build /api /api
ARG PORT
EXPOSE ${PORT}
USER nobody
ENTRYPOINT [ "/api" ]

FROM alpine:3.23 AS worker
WORKDIR /
COPY --from=build /worker /worker
USER nobody
ENTRYPOINT [ "/worker" ]