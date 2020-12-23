FROM golang:1.14 as builder

LABEL maintaner="Nick Mrozowski <nickmro@gmail.com>"

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/serve -i ./cmd/serve

FROM alpine:latest

LABEL maintainer="Nick Mrozowski <nickmro@gmail.com>"

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/bin/serve bin
COPY .env .

EXPOSE 3000

CMD ["./bin"]
