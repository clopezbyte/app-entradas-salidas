# Build Stage
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Image
FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/main /main
EXPOSE 8080
CMD ["/main"]
