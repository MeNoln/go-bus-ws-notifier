FROM golang:alpine AS build

WORKDIR /config
COPY ../config/* .

WORKDIR /app

COPY ../go.* .
RUN go mod download
RUN go mod verify

COPY worker .

ENV APP_ENV=dev

ENTRYPOINT [ "./worker" ]