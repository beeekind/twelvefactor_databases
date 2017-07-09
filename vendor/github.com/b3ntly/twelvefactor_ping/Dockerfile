# Use Docker 1.7 multi-stage builds to minimize Docker image sizes
# 1) From a basic image with Glide installed: build our src files into a binary
# 2) Copy that binary to an even smaller container representing the final image
FROM billyteves/alpine-golang-glide:1.2.0

# Copy all files in this directory into the requisite path
ADD . /go/src/github.com/b3ntly/twelvefactor_ping
WORKDIR /go/src/github.com/b3ntly/twelvefactor_ping

# Use Glide to install all of our dependencies
RUN glide install

# Build are binary for linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /root/

# Copy in our binary from the first stage
COPY --from=0 /go/src/github.com/b3ntly/twelvefactor_ping/main .

# Set the entry point to the binary
CMD ["./main"]