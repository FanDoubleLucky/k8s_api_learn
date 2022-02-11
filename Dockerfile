FROM golang:1.16-alpine

ENV GOPATH /go/
RUN go env -w GOPROXY="https://goproxy.cn/"
WORKDIR /go/src/k8s_api_learn
COPY . /go/src/k8s_api_learn
RUN go mod tidy
RUN go build .

EXPOSE 8111
ENTRYPOINT ["./k8s_api_learn"]