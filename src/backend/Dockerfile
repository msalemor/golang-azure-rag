# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app
COPY . .
COPY wwwroot wwwroot

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gorag

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/gorag gorag
COPY --from=build-stage /app/wwwroot wwwroot
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/gorag"]