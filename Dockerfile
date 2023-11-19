FROM --platform=linux/amd64 golang:1.21-alpine AS builder

WORKDIR /app

COPY . ./

RUN go mod vendor
RUN go mod tidy

RUN go build -o /web-scrapper

RUN ls
RUN ls /
RUN ls /app

FROM scratch

WORKDIR /app

COPY --from=builder /web-scrapper /app

EXPOSE 3000

CMD ["/app/web-scrapper"]
