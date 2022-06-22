FROM debian

RUN apt-get update && \
    apt-get install -y bluetooth git golang && \
    rm -rf /var/lib/apt/lists/*

RUN cd /tmp && \
    git clone https://github.com/go-delve/delve && \
    cd delve && \
    go build ./cmd/dlv/ && \
    mv dlv /usr/bin && \
    cd .. && \
    rm -rf delve

ENTRYPOINT while :; do :; done & kill -STOP $! && wait $!