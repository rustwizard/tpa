FROM golang:alpine as builder
RUN apk add --update git alpine-sdk

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


FROM alpine
RUN apk add --update --no-cache ca-certificates openssl
COPY --from=builder /app/tpa /app/
EXPOSE 8080
ENTRYPOINT ["/app/tpa"]