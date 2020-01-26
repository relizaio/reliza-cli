FROM relizaio/versioning as version-stage
ARG MODIFIER=Snapshot
ARG CI_ENV=noci
ARG GIT_COMMIT=git_commit_undefined
ARG GIT_BRANCH=git_branch_undefined
RUN echo "version=$(/usr/bin/java -jar /app/versioning.jar -s YYYY.0M.Modifier+Metadata -i $MODIFIER -m $CI_ENV)" > /tmp/version && echo "commit=$GIT_COMMIT" >> /tmp/version && echo "branch=$GIT_BRANCH" >> /tmp/version

FROM golang:1.13.6-buster as build-stage
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:stable-20191224-slim as release-stage
ARG CI_ENV=noci
ARG GIT_COMMIT=git_commit_undefined
ARG GIT_BRANCH=git_branch_undefined
RUN apt-get update
RUN apt-get -y install ca-certificates
RUN mkdir /app
RUN useradd apprunner && chown apprunner:apprunner /app
COPY --from=build-stage --chown=apprunner:apprunner /go/bin/app /app/app
COPY --from=version-stage --chown=apprunner:apprunner /tmp/version /app/version
USER apprunner
RUN mkdir /app/localdata
LABEL git_commit $GIT_COMMIT
LABEL git_branch $GIT_BRANCH
LABEL ci_environment $CI_ENV
ENTRYPOINT ["/app/app"]