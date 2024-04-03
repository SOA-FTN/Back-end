FROM golang:latest

RUN mkdir /build
WORKDIR /build

COPY . .


RUN go get -d -v
RUN cd /build && go build main.go

EXPOSE 8083

ENTRYPOINT ["/build/main"]