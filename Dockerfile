FROM alpine:3.17.3 as base-alpine
RUN apk add --no-cache ffmpeg

FROM golang:1.21.0 as assets-dependencies
COPY go.mod /go/src/project/
COPY go.sum /go/src/project/
WORKDIR /go/src/project
RUN go mod download
RUN go mod verify

FROM assets-dependencies as assets-build
COPY . /go/src/project/
WORKDIR /go/src/project
RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o assets cmd/application/main.go

FROM base-alpine
COPY --from=assets-build /go/src/project/assets /etc/assets/assets
COPY --from=assets-build /go/src/project/config /etc/assets/config/
WORKDIR /etc/assets
EXPOSE 8080
CMD ["./assets"]
