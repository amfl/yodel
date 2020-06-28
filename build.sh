#!/bin/sh

go build -v -o build/yodel ./cmd/yodel
# http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
# CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .
