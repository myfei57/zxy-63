FROM golang:1.23 AS builder
WORKDIR /src
COPY . .
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
RUN go build -mod=vendor -o /waterplant ./cmd/waterplant

FROM golang:1.23
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
COPY --from=builder /src /app
COPY --from=builder /waterplant /usr/local/bin/waterplant
CMD ["waterplant", "-addr", ":8080"]
