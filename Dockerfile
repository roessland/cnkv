# Stage 1: Compile the binary in a containerized Docker environment
#
FROM golang:1.16 as build

# Copy the source files from the host
COPY . /src

# Set the working directory to the same place we copied the code
WORKDIR /src

# Build the binary!
RUN CGO_ENABLED=0 GOOS=linux go build -o cnkv.bin ./cnkv

# Use a "scratch" image, which contains no distribution files
FROM scratch

# Copy the binary from the build container
COPY --from=build /src/cnkv.bin .

# Tell Docker we'll be using port 8080
EXPOSE 8080

# Tell Docker to run this command on a "docker run"
CMD ["/cnkv.bin"]