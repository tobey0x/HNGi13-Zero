FROM golang:1.25.3 AS builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy && go build -o server .


FROM gcr.io/distroless/base-debian11

USER nonroot

COPY --from=builder /app/server /server


EXPOSE 8080

CMD ["/server"]