# --- build stage ---
# Go 1.26: matches the `go 1.26.0` directive in go.mod (bumped when client-go
# was added). Build is fully offline via the committed-to-context vendor/ dir.
FROM golang:1.26 AS build
WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/apikit ./cmd/apikit

# --- runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache bash curl jq ca-certificates

COPY --from=build /out/apikit /usr/local/bin/apikit

ENTRYPOINT ["/bin/bash"]