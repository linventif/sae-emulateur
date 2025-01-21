FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN git reset --hard 53b2ca94937d17e62185f15835bad63f84bf331c

RUN go build -o sae-emulateur .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sae-emulateur .

ENTRYPOINT ["./sae-emulateur"]
