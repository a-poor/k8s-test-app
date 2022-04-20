FROM golang:1.18 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./build/testapp


FROM scratch
WORKDIR /app
COPY --from=build /app/build/testapp .

ENV API_ENV=dev

EXPOSE 80
ENTRYPOINT [ "/app/testapp" ]
