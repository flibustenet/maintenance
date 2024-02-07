#FROM amd64/golang:latest
FROM golang:1.22.0-alpine3.19
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o maintenance

#FROM gcr.io/distroless/static-debian12

# Copy the binary to the production image from the builder stage.
#COPY --from=builder /app/maintenance /app/maintenance
WORKDIR /app


# Run the web service on container startup.
ENTRYPOINT ["/app/maintenance"]

