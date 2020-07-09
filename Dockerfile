ARG goVersion=1.14.3-alpine

FROM golang:${goVersion} as builder
ARG gitCommit=""
ARG buildStamp=""

ENV GO111MODULE=on

WORKDIR /app

ADD . .

RUN go build -ldflags "-s -w -X 'main.gitCommit=${gitCommit}' -X 'main.buildStamp=${buildStamp}'" -o sidecar

FROM alpine

LABEL author=renjie email=penitence.rj@gmail.com

ENV app=sidecar

COPY --from=builder /app/sidecar /usr/local/bin/

CMD ["sidecar"]