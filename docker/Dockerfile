FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \ 
    go build -ldflags "-s -w" -o /app/BcvScraper .

FROM alpine:latest AS runner
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/BcvScraper .
USER 1000:1000
CMD ["./BcvScraper"]