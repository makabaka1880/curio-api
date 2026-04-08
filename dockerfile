FROM golang:latest AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app . && ls -la app

FROM alpine:latest AS run
WORKDIR /app
COPY --from=build /build/app .
ENTRYPOINT ["/app/app"]