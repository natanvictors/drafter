FROM golang:1.26.4

WORKDIR /drafter

COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -o /docker-drafter ./cmd/drafter


CMD [ "/docker-drafter" ]
