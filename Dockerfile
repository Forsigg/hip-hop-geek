FROM golang:1.22.2-bookworm

WORKDIR /app

COPY go.mod go.sum ./
COPY ./.env ./.env
RUN go mod download

COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./cmd ./cmd

RUN CGO_ENABLED=1 GOOS=linux go build -o /hip-hop-geek-bot ./cmd/app/main.go


CMD ["/hip-hop-geek-bot"]
