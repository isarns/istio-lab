FROM golang:1.22 AS build-stage

ARG APP_DIR
ARG PORT

WORKDIR /opt

COPY utils ../utils
COPY ${APP_DIR} ./

RUN go mod download

COPY ${APP_DIR}/*.go ./


RUN CGO_ENABLED=0 GOOS=linux go build -o /app

# Run the tests in the container
FROM build-stage AS run-test-stage

RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

ARG PORT
ENV PORT=${PORT}

WORKDIR /

COPY --from=build-stage /app /app

EXPOSE ${PORT}

USER nonroot:nonroot

ENTRYPOINT ["/app"]