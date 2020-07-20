FROM golang:1.14-alpine as builder
# Force the go compiler to use modules
ENV GO111MODULE=on
# Create the user and group files to run unprivileged 
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk update && apk add --no-cache git ca-certificates tzdata
RUN mkdir /build
# COPY . /build/
WORKDIR /build
# Get the deps - Do this separately for ease of iteration
COPY go.mod go.sum /build/
RUN go mod download
# Import the code
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o yodel ./cmd/yodel

##############################

FROM scratch AS final
# Import the time zone files
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Import the user and group files
COPY --from=builder /user/group /user/passwd /etc/
# Import the CA certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled go executable
COPY --from=builder /build/yodel /
WORKDIR /
# Run as unprivileged
USER nobody:nobody
ENV YODEL_CONFIG_PATH=/config
VOLUME /config
ENTRYPOINT ["/yodel"]
