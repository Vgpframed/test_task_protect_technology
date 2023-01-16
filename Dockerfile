
FROM golang:alpine as BUILDER

RUN apk update && apk add --no-cache git ca-certificates
COPY . /app
RUN cd /app \
    && go mod download \
    && go generate \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
    -o /app/app.o ./cmd/app/

FROM scratch
COPY --from=BUILDER /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=BUILDER /etc/passwd /etc/passwd
COPY --from=BUILDER /app/app.o /usr/bin/app
ENTRYPOINT ["app"]
