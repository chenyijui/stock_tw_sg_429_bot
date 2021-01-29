FROM golang:1.13-alpine
WORKDIR /go/src/stock_tw_sg_429_bot
ADD . /go/src/stock_tw_sg_429_bot
RUN cd /go/src/stock_tw_sg_429_bot && go build
RUN go build -o app
EXPOSE 8787
ENTRYPOINT ./app