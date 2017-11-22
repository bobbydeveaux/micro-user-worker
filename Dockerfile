FROM golang:1.8
COPY       micro-user-worker /bin/micro-user-worker
ENTRYPOINT ["/bin/micro-user-worker"]
