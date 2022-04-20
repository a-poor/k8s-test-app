FROM golang:1.18-bullseye AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./build/testapp


FROM scratch
WORKDIR /app
COPY --from=build /app/build/testapp .

ENV API_ENV="dev"
ENV API_PORT="1323"
ENV API_HTTPS=""

# USER nonroot:nonroot

EXPOSE 80
ENTRYPOINT [ "/app/testapp" ]
