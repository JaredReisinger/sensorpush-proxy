FROM golang:1.19.5-alpine AS builder

RUN \
  apk update && \
  apk add --no-cache \
    ca-certificates \
    git \
    tzdata \
    zip \
  && \
  update-ca-certificates && \
  zip -q -r -0 /zoneinfo.zip /usr/share/zoneinfo && \
  wget -O - https://taskfile.dev/install.sh | sh -s -- -d -b /usr/local/bin

WORKDIR /app

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

RUN task build

# ----------------------------------------------------------------------
# now create a minimal docker image
FROM scratch

ARG SOURCE_COMMIT
ARG DOCKER_TAG
ARG BUILD_DATE

LABEL \
  org.opencontainers.image.title="sensorpush-proxy" \
  org.opencontainers.image.description="" \
  org.opencontainers.image.licenses="" \
  org.opencontainers.image.url="https://github.com/JaredReisinger/sensorpush-proxy" \
  org.opencontainers.image.documentation="https://github.com/JaredReisinger/sensorpush-proxy#readme" \
  org.opencontainers.image.source="https://github.com/JaredReisinger/sensorpush-proxy" \
  org.opencontainers.image.revision="${SOURCE_COMMIT}" \
  org.opencontainers.image.version="${DOCKER_TAG}" \
  org.opencontainers.image.created="${BUILD_DATE}" \
  org.label-schema.schema-version="1.0" \
  org.label-schema.name="sensorpush-proxy" \
  org.label-schema.description="" \
  org.label-schema.license="" \
  org.label-schema.url="https://github.com/JaredReisinger/sensorpush-proxy" \
  org.label-schema.vcs-url="https://github.com/JaredReisinger/sensorpush-proxy" \
  org.label-schema.vcs-ref="${SOURCE_COMMIT}" \
  org.label-schema.version="${DOCKER_TAG}" \
  org.label-schema.build-date="${BUILD_DATE}"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# do we need TZ info?
# COPY --from=builder /zoneinfo.zip /
# ENV ZONEINFO /zoneinfo.zip

COPY --from=builder /app/build/sensorpush-proxy /sensorpush-proxy

#EXPOSE 80
EXPOSE 5375

# don't run as root?
USER 1000:1000
ENTRYPOINT ["/sensorpush-proxy"]
