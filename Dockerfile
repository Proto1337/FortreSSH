# This Dockerfile uses two stages.
# The first stage builds the binary.
# The second stage executes the binary on default port 2222.
# It is recommended to keep it on 2222 and use a DNAT.

# Stage 1: Building
# Use the Go Alpine image as building environment
FROM docker.io/library/golang:alpine as builder
WORKDIR /pit
COPY fortressh.go /pit
RUN go build fortressh.go

# Stage 2: Running
# Run in minimal busybox environment
FROM docker.io/library/busybox:musl
ENV LANG=C.UTF-8
WORKDIR /pit
# Always use an unpriviliged user for your Docker stuff!
RUN adduser -D appuser
USER appuser
COPY --from=builder --chown=appuser:appuser /pit/fortressh /pit
# Expose the port
EXPOSE 2222/tcp
ENTRYPOINT ["/pit/fortressh"]
