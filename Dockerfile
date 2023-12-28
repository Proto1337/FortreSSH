# This Dockerfile uses two stages.
# The first stage builds the binary.
# The second stage executes the binary on default port 2222.
# It is recommended to keep the port on 2222 and use a DNAT.

# Stage 1: Building
# Use the Go Alpine image as building environment
FROM docker.io/library/golang:alpine as builder
# Create passwd here. scratch can not.
RUN adduser -HDs /bin/false pitter
WORKDIR /pit
RUN apk add llvm
COPY fortressh.go /pit
RUN go build fortressh.go && llvm-strip fortressh

# Stage 2: Running
# Run in minimal environment
FROM scratch
ENV LANG=C.UTF-8
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder --chown=pitter /pit/fortressh /fortressh
# Always use an unprivileged user
USER pitter
# Expose the port
EXPOSE 2222/tcp
# Start the tarpit
ENTRYPOINT ["/fortressh"]
