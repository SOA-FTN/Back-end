FROM golang:latest

RUN mkdir /build
WORKDIR /build

COPY ./. .
COPY ./main.go .
COPY ./go.mod .
COPY ./go.sum .
COPY ./DtoObjects ./DtoObjects
COPY ./handler ./handler
COPY ./model ./model
COPY ./repo ./repo
COPY ./service ./service
COPY ./monitoring ./monitoring


RUN go get -d -v
RUN cd /build && go build main.go

EXPOSE 8082

ENTRYPOINT ["/build/main"]