FROM golang:alpine as builder

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app cmd/main.go

FROM gcr.io/distroless/static-debian11

EXPOSE 8192

COPY --from=builder /go/bin/app /app

CMD ["/app"]
