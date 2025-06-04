FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.sum ./
COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o app cmd/api/*.go

# -----------------------------------------------------
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/.env.production .

EXPOSE 8090
CMD [ "./app" ]