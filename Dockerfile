FROM golang:1.19-alpine AS build
LABEL MAINTAINER github.com/arizon-dread

WORKDIR /usr/local/go/src/github.com/arizon-dread/status-checker-api
COPY businesslayer ./businesslayer
COPY models ./models
COPY main.go go.mod go.sum ./

RUN apk update && apk add --no-cache git
RUN go build -v -o /usr/local/bin/status-checker-api/ ./...


FROM golang:1.19-alpine AS final
WORKDIR /go/bin
ENV GIN_MODE=release
RUN apk add --no-cache libc6-compat musl-dev
COPY --from=build /usr/local/bin/status-checker-api/ /go/bin/
EXPOSE 8080

ENTRYPOINT [ "./status-checker-api" ]
