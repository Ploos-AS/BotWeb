FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN go test ./... && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/botweb ./cmd/botweb
FROM alpine:3.22
RUN adduser -D -H botweb
COPY --from=build /out/botweb /usr/local/bin/botweb
USER botweb
ENTRYPOINT ["botweb"]
