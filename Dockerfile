FROM golang:1.19-alpine as builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN git clone https://github.com/FeranyDev/mini-push.git && cd mini-push && go install && go build -o /app/bin/mini-push


FROM alpine:latest

WORKDIR /app

MAINTAINER feranydev@gmail.com

ENV MINI_PUSH_CONFIG /app/config.yaml

COPY --from=builder /app/mini-push/config.yaml /app/config.yaml
COPY --from=builder /app/bin/mini-push /app/mini-push

EXPOSE 3000/tcp

CMD ["./mini-push"]