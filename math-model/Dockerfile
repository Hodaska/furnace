FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o math-model .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/math-model .
EXPOSE 4840
CMD ["./math-model"]
