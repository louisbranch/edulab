FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/server ./cmd/server

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/web/templates /var/www/templates
COPY --from=build /go/src/app/web/assets /var/www/assets
EXPOSE 80
ENTRYPOINT /go/bin/server --port 80