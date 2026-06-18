# Runtime-only image for the in-cluster apikit pod.
# It packages a PREBUILT linux/amd64 binary at bin/apikit in the build context
# (produced by the `build-and-test` job and downloaded as an artifact).
# No Go build happens here — the binary is compiled once in CI, then shipped.
FROM alpine:3.20

RUN apk add --no-cache bash curl jq ca-certificates

COPY bin/apikit /usr/local/bin/apikit
RUN chmod +x /usr/local/bin/apikit

ENTRYPOINT ["/bin/bash"]
