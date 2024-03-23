FROM golang:1.22.0-alpine as build

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w" -o godo cmd/godo/main.go 

FROM alpine:latest
COPY --from=build /build/godo /godo
COPY --from=build /build/static /static
EXPOSE 8080

CMD ["/godo"]