FROM golang:1.20

WORKDIR /treechat

COPY go.mod .

COPY . .

RUN make build

CMD ["make", "run"]
