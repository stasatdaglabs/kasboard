FROM golang:1.16-alpine AS build

# gcc and libc (musl-dev) are required by kaspad
# git is required to build kaspad by commit-hash
RUN apk add gcc musl-dev git

RUN mkdir /build
WORKDIR /build

COPY processing/go.mod .
COPY processing/go.sum .

ARG KASPAD_VERSION

# Exit if KASPAD_VERSION is not set
RUN if [ -z "${KASPAD_VERSION}" ]; then exit 1; fi

# Replace the kaspad version in go.mod
ARG KASPAD_REPOSITORY='github.com/kaspanet/kaspad'
RUN go mod edit -replace github.com/kaspanet/kaspad=$KASPAD_REPOSITORY@$KASPAD_VERSION

RUN go mod download

RUN cp go.mod go.mod.bu
RUN cp go.sum go.sum.bu

COPY processing .

# Restore go.mod and go.sum because `COPY processing .` overwrote them
RUN mv go.mod.bu go.mod
RUN mv go.sum.bu go.sum
RUN go mod tidy

RUN GOOS=linux go build -o processing .

FROM alpine
WORKDIR /app
COPY --from=build /build/processing /app/

RUN mkdir -p /app/database/migrations
COPY --from=build /build/database/migrations/ /app/database/migrations/
