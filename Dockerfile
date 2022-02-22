FROM golang:1.17 as builder

RUN mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app

RUN go mod tidy
RUN CGO_ENABLED=0 go build

FROM alpine

COPY --from=0 /go/src/app/NTLM-proxy .
ENTRYPOINT ["./NTLM-proxy"]

EXPOSE 8080