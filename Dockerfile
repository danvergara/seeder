FROM golang:1.16.6-buster AS builder

WORKDIR /src/app

# install system dependencies
RUN apt-get update \
  && apt-get -y install netcat curl \
  && apt-get clean

COPY go.* Makefile ./
RUN go mod download

COPY scripts ./scripts
RUN make install-migrate

COPY . . 

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o seeder ./cli

FROM alpine:3.14 AS bin

RUN apk add --no-cache ca-certificates git make musl-dev go

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
WORKDIR /seeder

LABEL org.opencontainers.image.documentation="https://github.com/danvergara/seeder" \
	org.opencontainers.image.source="https://github.com/danvergara/seeder" \
	org.opencontainers.image.title="seeder"

COPY --from=builder /src/app/seeder /usr/local/bin/seeder
RUN ln -s /usr/local/bin/seeder /seeder

ENTRYPOINT [ "seeder" ]
