FROM golang:1.19-alpine AS build

WORKDIR /app/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /app/auditor

FROM alpine

COPY --from=build /app/auditor /app/auditor

ENTRYPOINT ["/app/auditor"]