FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY src ./

RUN go build -o main .

FROM alpine:3.17
COPY --from=build /app/main /main
ENTRYPOINT [ "/main" ]
