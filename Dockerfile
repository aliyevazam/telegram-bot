FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o main cmd/main.go

EXPOSE 8000

CMD ["/app/main"]