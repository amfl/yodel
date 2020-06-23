FROM golang:1.14.3-buster

WORKDIR /src

# Grab dependencies
COPY src/go.mod src/go.sum /src/
RUN go mod download

# Compile code
ENV VERSION=0.1.0 \
    YODEL_CONFIG_PATH=/config

COPY src /src
# RUN go install \
#     -installsuffix "static" \
#     -ldflags "-X $(go list -m)/pkg/version.Version=${VERSION}"

VOLUME /config
ENTRYPOINT /go/bin/ldap-lookup

# Set cache dir to somewhere we can write to
#   go env GOCACHE
# ENV GOCACHE=/tmp/workdir/.cache/go-build
