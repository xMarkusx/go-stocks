FROM golang:1.20

WORKDIR /go/src/stock-monitor

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /stock-monitor

EXPOSE 8080

CMD ["/stock-monitor"]
