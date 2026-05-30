# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.26.2 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build .

# Deploy the application bin0 ary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/osmtools /osmtools

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/osmtools"]