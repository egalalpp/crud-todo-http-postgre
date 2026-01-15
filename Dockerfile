FROM golang:1.25.5-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -o /app/exe ./cmd/main.go


CMD ["/app/exe"]



