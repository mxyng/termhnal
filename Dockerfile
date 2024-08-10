FROM golang:alpine

WORKDIR /go/src/github.com/mxyng/termhnal

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o termhnal

FROM alpine
ENV TERM=xterm-256color
COPY --from=0 /go/src/github.com/mxyng/termhnal/termhnal /usr/bin/termhnal
ENTRYPOINT ["/usr/bin/termhnal"]
