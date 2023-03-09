FROM golang:1.20

WORKDIR /treechat

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN make build

CMD ["make", "run"]
