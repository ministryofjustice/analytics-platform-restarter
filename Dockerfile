# See: https://www.cloudreach.com/blog/containerize-this-golang-dockerfiles/

# Builder image stage
FROM golang:1.12-alpine AS builder

RUN apk update \
      && apk add --no-cache \
      ca-certificates \
      git \
      make

WORKDIR /go/src/github.com/ministryofjustice/analytics-platform-restarter

COPY vendor/ vendor/
COPY Makefile ./
COPY *.go ./
COPY go.mod ./
COPY go.sum ./

RUN go mod verify
RUN make test
RUN make vet

# NOTE: statically compiled as final image is based on "scratch"
RUN make static

# Binary image stage
FROM scratch
WORKDIR /bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ministryofjustice/analytics-platform-restarter/restarter .

CMD ["/bin/restarter"]
