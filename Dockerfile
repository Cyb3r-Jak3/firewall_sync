FROM golang:1.18.0-alpine as build

WORKDIR /go/src/app
COPY . /go/src/app

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go build -buildvcs=false -o /go/bin/app ./...

FROM gcr.io/distroless/static as distroless
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
CMD ["-c", "config.yml"]

FROM alpine:latest as alpine
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
CMD ["-c", "config.yml"]
