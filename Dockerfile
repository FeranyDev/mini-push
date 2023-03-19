FROM golang:1.19-alpine as builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

RUN git clone https://github.com/FeranyDev/mini-push.git && cd mini-push && go mod && go build -o /app/bin/mini-push


FROM alpine:latest

COPY --from=builder /app/bin/mini-push /app/bin/mini-push

CMD ["/app/bin/mini-push"]