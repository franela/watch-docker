FROM golang:1.8

# Copy the code inside the container and set the workdir
COPY ./ /go/src/github.com/franela/watch-docker
WORKDIR /go/src/github.com/franela/watch-docker

# Build the binary
RUN go get ./...
CMD ["watch-docker"]
