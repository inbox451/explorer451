FROM gcr.io/distroless/static-debian12:nonroot

# Default environment variables
ENV EXPLORER451_SERVER_HTTP_PORT=":8080" \
  EXPLORER451_LOGGING_LEVEL="info" \
  EXPLORER451_LOGGING_FORMAT="json"

# Copy binary
COPY explorer451 /explorer451

# Expose ports
EXPOSE 8080

# Run the application
ENTRYPOINT ["./explorer451"]
