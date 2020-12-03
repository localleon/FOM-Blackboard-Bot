FROM golang:alpine

# Make App dir
RUN mkdir /app
ADD ./*go* /app

# Build Application
WORKDIR /app
RUN go build -o ocBot .

# Set entrypoint
CMD ["/app/ocBot"]