FROM golang:1.15-alpine AS build

RUN mkdir /build
WORKDIR /build

COPY processing/go.mod .
COPY processing/go.sum .

RUN go mod download

COPY processing .

RUN GOOS=linux go build -o processing .

FROM alpine
WORKDIR /app
COPY --from=build /build/processing /app/

RUN mkdir -p /app/database/migrations
COPY --from=build /build/database/migrations/ /app/database/migrations/
