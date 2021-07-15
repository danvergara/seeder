FROM golang:1.16.6-buster AS builder

WORKDIR /src/app

# install system dependencies
RUN apt-get update \
  && apt-get -y install netcat curl \
  && apt-get clean

COPY go.* ./
RUN go mod download
COPY . . 

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH}  go build -o seeder .

RUN make install-migrate

FROM scratch AS bin

LABEL org.opencontainers.image.documentation="https://github.com/danvergara/seeder" \
	org.opencontainers.image.source="https://github.com/danvergara/seeder" \
	org.opencontainers.image.title="seeder"

COPY --from=builder /src/app/seeder /bin/seeder
