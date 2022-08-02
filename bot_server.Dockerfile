FROM golang:1.18.4-alpine

WORKDIR /app

RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

COPY ./go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
COPY ./entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "entrypoint.sh"]