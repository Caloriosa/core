FROM golang:1.10-alpine

ENV PATH=$PATH:/go/src/core
ENV HOME=/tmp

RUN mkdir -p /config

WORKDIR /go/src/core

COPY . .
RUN  ./build.sh; \
     chmod -R 777 /tmp/.cache

ENTRYPOINT ["bin/caloriosa-server"]
CMD ["-config", "/config/config.yaml", "-logtostderr"]
