FROM golang:1.17.0-alpine as build

WORKDIR /go/src/app
COPY . /go/src/app

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go build -o /go/bin/app

FROM gcr.io/distroless/static
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
CMD ["-c", "config.yml"]