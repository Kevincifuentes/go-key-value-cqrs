# syntax=docker/dockerfile:1

# Stage 1: Build the application
FROM golang:1.23.4 AS build-stage

WORKDIR /app

# Copy dependency files and download modules
COPY . .
RUN go mod download

# Copy the source code and build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/keyvalueserver .

# Stage 2: Run tests
FROM build-stage AS test-stage
RUN bash ./runAllTests.sh

# Stage 3: Export test results
FROM scratch AS test-out
COPY --from=test-stage /app/assets/ /

# Stage 4: Create a minimal runtime image
FROM alpine:latest AS runtime-stage
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build-stage /app/keyvalueserver /app/api/keyvalue/api.yml /app/.env* ./

# Default port
ARG SERVER_PORT=8080
EXPOSE $SERVER_PORT
# Debug mode
ARG DEBUG_PORT=8081
EXPOSE $DEBUG_PORT

# Set the command to run the server
ENV SERVER_PORT=$SERVER_PORT
ENV DEBUG_SERVER_PORT=$DEBUG_PORT
ENV OPENAPI_RELATIVE_PATH="./api.yml"
CMD ["./keyvalueserver"]