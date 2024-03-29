FROM golang:1.19 as builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o maintenance

FROM gcr.io/distroless/static-debian11

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/maintenance /app/maintenance
WORKDIR /app


# Run the web service on container startup.
CMD ["/app/maintenance"]

