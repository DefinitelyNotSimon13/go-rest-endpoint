FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o /tmp/test-endpoint ./api/main.go

FROM alpine:edge

COPY --from=build /tmp/test-endpoint /sbin/test-endpoint

CMD /sbin/test-endpoint
