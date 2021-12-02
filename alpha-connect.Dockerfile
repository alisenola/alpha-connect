FROM golang:1.16 as builder

ENV GO111MODULE=on
ENV GOPRIVATE="gitlab.com/lmeunier,gitlab.com/alphaticks,gitlab.com/alphaticks"
RUN git config --global url."https://lmeunier:PQDWxcmr6p3fyvpVHrAy@gitlab.com/".insteadOf "https://gitlab.com/"

WORKDIR /src

COPY ./go.mod ./go.sum ./
ENV GOPROXY=http://127.0.0.1:8050
RUN go mod download

FROM builder as app_builder

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM scratch

COPY --from=app_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app_builder /etc/passwd /etc/passwd
COPY --from=app_builder /src/alpha-connect /alpha-connect

ENTRYPOINT ["/alpha-archiver"]
