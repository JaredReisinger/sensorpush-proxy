# This is a goreleaser-style Dockerfile, where the binary is built by
# goreleaser, *not* in a Docker build stage!

FROM alpine AS certs

RUN \
  apk update && \
  apk add --no-cache ca-certificates && \
  update-ca-certificates

# ---------------------------------------------------------------------------
FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 5375
USER 1000:1000
ENTRYPOINT [ "/sensorpush-proxy" ]

COPY sensorpush-proxy /
