FROM golang:1.22 as builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o maintenance

FROM gcr.io/distroless/static-debian12


COPY --from=builder /app/maintenance /app/maintenance
WORKDIR /app

ENTRYPOINT ["/app/maintenance"]

