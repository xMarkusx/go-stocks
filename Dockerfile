FROM golang:1.20

WORKDIR /go/src/stock-monitor

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go env

RUN CGO_ENABLED=0 GOOS=linux go build -o /stock-monitor

EXPOSE 8080

CMD ["/stock-monitor"]
