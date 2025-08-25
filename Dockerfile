# Use a minimal base image (Alpine for lightweight)
FROM alpine:3.18

# Install ca-certificates for HTTPS (if needed by CLI)
RUN apk add --no-cache ca-certificates

# Copy the pre-built threadbolt binary
COPY dist/threadbolt-linux-amd64 /usr/local/bin/threadbolt

# Set execute permissions
RUN chmod +x /usr/local/bin/threadbolt

# Set the entrypoint to the threadbolt CLI
ENTRYPOINT ["threadbolt"]

# Default command (can be overridden)
CMD ["--help"]
