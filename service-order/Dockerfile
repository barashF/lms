FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /var/log/service-order && chmod 755 /var/log/service-order
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/service /service-order
COPY --from=builder /app/.env /
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/service-order"]
