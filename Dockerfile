# This Dockerfile uses two stages.
# The first stage builds the binary.
# The second stage executes the binary on port 2222.

# Stage 1: Building
FROM golang:alpine as builder
WORKDIR /pit
ADD fortressh.go .
RUN go build fortressh.go

# Stage 2: Running
FROM alpine:latest
WORKDIR /pit
RUN passwd -l root
# Always use an unpriviliged user for your Docker stuff!
RUN adduser -D appuser
USER appuser
COPY --from=builder --chown=appuser:appuser /pit/fortressh /pit
EXPOSE 2222/tcp
ENTRYPOINT ["/pit/fortressh"]
