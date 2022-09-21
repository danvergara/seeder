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
