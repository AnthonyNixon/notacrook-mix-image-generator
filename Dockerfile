FROM golang:1.14-alpine as builder

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/anthonynixon/notacrook-mix-image-generator
COPY go.mod .
RUN go mod tidy

FROM builder as binary_builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /notacrook-mix-image-generator .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=binary_builder /notacrook-mix-image-generator /notacrook-mix-image-generator
COPY not_a_crook_logo.png /not_a_crook_logo.png
COPY impact.ttf /impact.ttf

CMD ["./notacrook-mix-image-generator"]