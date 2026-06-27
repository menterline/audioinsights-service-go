FROM golang:1.26.4 AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/audioinsights-service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/audioinsights-service /audioinsights-service
EXPOSE 8080
ENTRYPOINT ["/audioinsights-service"]
