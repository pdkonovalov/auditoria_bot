FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /auditoria_bot

FROM gcr.io/distroless/static-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /auditoria_bot /auditoria_bot

USER nonroot:nonroot

ENTRYPOINT ["/auditoria_bot"]