FROM golang:1.21-alpine AS builder

WORKDIR /sso_project

RUN apk --no-cache add bash git gcc gettext musl-dev sqlite-dev

COPY ["app/go.mod","app/go.sum","./"]
RUN go mod download

COPY app ./
RUN go build -o ./app cmd/sso/main.go

FROM alpine AS runner

COPY --from=builder /sso_project /

ENTRYPOINT ["/app"]
CMD ["--config=./config/local.yaml"]