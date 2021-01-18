FROM alpine:3.11

ARG BINARY={{ .App.ProjectName }}
ENV BINARY ${BINARY}

RUN adduser -D -s /bin/sh ragnaros
WORKDIR /home/ragnaros

COPY ${BINARY} /home/ragnaros
COPY resources /home/ragnaros/resources

# alpine tricky for golang
RUN apk update && \
    apk add -u dumb-init && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8099

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/bin/sh", "-c", "/home/ragnaros/${BINARY}"]