FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy 

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app

FROM ubuntu:latest

WORKDIR /root/
COPY --from=builder /my_app .
COPY ./web ./web 
COPY .env . 

CMD ["./my_app"]