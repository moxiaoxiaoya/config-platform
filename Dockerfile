FROM golang:1.20
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-config-center

EXPOSE 8080

CMD ["/docker-config-center"]