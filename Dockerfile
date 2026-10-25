FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /daily .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /daily /daily
USER nonroot:nonroot
ENTRYPOINT ["/daily"]
