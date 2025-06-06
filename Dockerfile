FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.sum ./
COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o app cmd/api/main.go

# -----------------------------------------------------
FROM alpine:3.22
COPY --from=builder /app/app .
COPY --from=builder /app/.env.production .

RUN mkdir -p migrations
COPY --from=builder /app/internal/db/migration/* /migrations

RUN mkdir -p uploads

EXPOSE 8099
CMD [ "./app" ]