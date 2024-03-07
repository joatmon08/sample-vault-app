FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY app.env ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /sample-vault-app

EXPOSE 8080

CMD ["/sample-vault-app"]