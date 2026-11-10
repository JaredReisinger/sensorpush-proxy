# This is a goreleaser-style Dockerfile, where the binary is built by
# goreleaser, *not* in a Docker build stage!
FROM alpine:latest AS certs

RUN \
  apk update && \
  apk add --no-cache ca-certificates && \
  update-ca-certificates

# ---------------------------------------------------------------------------
FROM scratch

ARG TARGETPLATFORM

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ${TARGETPLATFORM}/sensorpush-proxy /

EXPOSE 5375
USER 1000:1000

ENTRYPOINT [ "/sensorpush-proxy" ]
