FROM golang:1.8

# Install gin, for an auto-reloading server
RUN go get github.com/codegangsta/gin

# Copy the code inside the container and set the workdir
COPY ./ /go/src/github.com/franela/watch-docker
WORKDIR /go/src/github.com/franela/watch-docker

# Build the binary
RUN go get ./...
CMD ["watch-docker"]
