FROM golang
MAINTAINER travis.simon@nicta.com.au

# Copy across our src files
ADD . /go/src/github.com/travissimon/microservices/strlen

# Build server
RUN go install github.com/travissimon/microservices/strlen

WORKDIR /go/bin

CMD strlen

# Listen on 8080
EXPOSE 8080