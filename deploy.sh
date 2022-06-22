#!/bin/sh

CONTAINER_NAME="nuki-test-env"

CGO_ENABLED=0 go build -a -installsuffix cgo -gcflags="all=-N -l" -o go-nuki ./cmd/example/
docker cp go-nuki ${CONTAINER_NAME}:/tmp/go-nuki
docker exec -it ${CONTAINER_NAME} dlv --listen=:2345 --headless=true --api-version=2 exec /tmp/go-nuki $@