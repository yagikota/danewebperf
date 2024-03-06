FROM golang:alpine AS builder
RUN apk add --no-cache unbound-dev build-base
COPY . /dane
WORKDIR /dane/cmd/letsdane
RUN go build -tags unbound

FROM alpine:latest
RUN apk add --no-cache unbound-libs
COPY --from=builder /dane /dane
WORKDIR /dane/cmd/letsdane
EXPOSE 8080

COPY ca/cert.crt ca/cert.key /root/.letsdane/
# Install tcpdump
RUN apk update && apk add tcpdump
# Create a directory for captured files
RUN mkdir /captured
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
