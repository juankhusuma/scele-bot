FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/server -ldflags="-s -w"

FROM scratch as app
COPY --from=builder /bin/server ./server
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env .env

ENTRYPOINT [ "./server" ]