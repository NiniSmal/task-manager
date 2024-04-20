FROM golang:1.21 AS builder
WORKDIR /app
ENV GOPRIVATE=gitlab.com/nina8884807/mail
RUN echo "machine gitlab.com login ninamusatova90 password glpat-61DSd-F9qJ4H9sZnqwwp" > $HOME/.netrc
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/main ./main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8021
ENTRYPOINT ["/app/main"]
