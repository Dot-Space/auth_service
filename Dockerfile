FROM golang:1.23.4-bookworm

WORKDIR /code

ENV GO111MODULE=on

COPY . .
EXPOSE 8080

RUN go install github.com/air-verse/air@latest

CMD ["air"]