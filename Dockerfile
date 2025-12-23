# syntax=docker/dockerfile:1

ARG PORT=8080

FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /api cmd/api/main.go

FROM alpine:3.23 AS api
WORKDIR /
COPY --from=build /api /api
ARG PORT
EXPOSE ${PORT}
USER nobody
ENTRYPOINT [ "/api" ]
