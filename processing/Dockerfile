FROM golang:1.16-alpine AS build

# gcc and libc (musl-dev) are required by kaspad
# git is required to build kaspad by commit-hash
RUN apk add gcc musl-dev git

RUN mkdir /build
WORKDIR /build

COPY processing/go.mod .
COPY processing/go.sum .

ARG KASPAD_COMMIT_HASH
RUN if [ -n "${KASPAD_COMMIT_HASH}" ]; then \
        # use replace instead of require - to propagate into nested dependancies
        go mod edit -replace github.com/kaspanet/kaspad=github.com/kaspanet/kaspad@$KASPAD_COMMIT_HASH ;\
    fi

RUN go mod download

RUN cp go.mod go.mod.bu
RUN cp go.sum go.sum.bu

COPY processing .

# Restore go.mod and go.sum because `COPY processing .` overwrote them
RUN mv go.mod.bu go.mod
RUN mv go.sum.bu go.sum

RUN GOOS=linux go build -o processing .

FROM alpine
WORKDIR /app
COPY --from=build /build/processing /app/

RUN mkdir -p /app/database/migrations
COPY --from=build /build/database/migrations/ /app/database/migrations/
