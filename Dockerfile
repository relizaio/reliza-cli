FROM golang:alpine3.12@sha256:6042b9cfb4eb303f3bdcbfeaba79b45130d170939318de85ac5b9508cb6f0f7e as build-stage
WORKDIR /go/src/app
RUN apk add git
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:3.12@sha256:185518070891758909c9f839cf4ca393ee977ac378609f700f60a771a2dfe321 as release-stage
ARG CI_ENV=noci
ARG GIT_COMMIT=git_commit_undefined
ARG GIT_BRANCH=git_branch_undefined
ARG VERSION=not_versioned
RUN mkdir /app
RUN adduser -u 1001 -D apprunner && chown apprunner:apprunner /app
COPY --from=build-stage --chown=apprunner:apprunner /go/bin/app /app/app
RUN mkdir /indir && chown apprunner:apprunner -R /indir
RUN mkdir /outdir && chown apprunner:apprunner -R /outdir
USER apprunner
RUN echo "version=$VERSION" > /app/version && echo "commit=$GIT_COMMIT" >> /app/version && echo "branch=$GIT_BRANCH" >> /app/version
RUN mkdir /app/localdata
LABEL git_commit $GIT_COMMIT
LABEL git_branch $GIT_BRANCH
LABEL ci_environment $CI_ENV
LABEL version $VERSION
ENTRYPOINT ["/app/app"]