FROM golang:1.19-alpine AS builder

WORKDIR /src

COPY . ./

RUN apk add build-base && go build  ./cmd/main.go

FROM alpine
WORKDIR /src
COPY --from=builder /src .

LABEL maintainers = "Maxat Kenzhalinov"
LABEL version = "1.0"

EXPOSE 8081

CMD ["/src/main"]