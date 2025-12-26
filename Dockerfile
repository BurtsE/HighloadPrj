FROM golang:1.25.4-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main/main.go

FROM golang:1.25.4-alpine
RUN apk --no-cache add ca-certificates
RUN addgroup -g 65532 nonroot && adduser -D -u 65532 nonroot -G nonroot
COPY --from=builder /app/main .
USER nonroot
EXPOSE 8080
CMD ["./main"]