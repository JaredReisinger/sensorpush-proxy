# This is a goreleaser-style Dockerfile, where the binary is built by
# goreleaser, *not* in a Docker build stage!
FROM scratch
ENTRYPOINT [ "/sensorpush-proxy" ]
COPY sensorpush-proxy /
