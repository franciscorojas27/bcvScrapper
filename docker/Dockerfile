FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY BcvScraper .
CMD ["./BcvScraper"]