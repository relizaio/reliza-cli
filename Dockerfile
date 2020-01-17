FROM golang:1.13.6-buster as build-stage
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
# ENTRYPOINT ["app"]

FROM debian:stable-20191224-slim as release-stage
RUN apt-get update
RUN apt-get -y install ca-certificates
RUN mkdir /app
RUN useradd apprunner && chown apprunner:apprunner /app
COPY --from=build-stage --chown=apprunner:apprunner /go/bin/app /app/app
USER apprunner
ENTRYPOINT ["/app/app"]