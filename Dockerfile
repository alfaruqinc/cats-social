FROM golang:1.22.2-alpine AS build

WORKDIR /app
COPY go.mod go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o /go/bin/main cmd/api/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/main /
CMD ["/main"]
