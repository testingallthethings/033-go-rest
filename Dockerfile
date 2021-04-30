# BUILD
FROM golang:latest as build

WORKDIR /service
ADD . /service

RUN chmod +x bin/entrypoint.sh
RUN apt update -yq
RUN apt install -y postgresql-client

RUN cd /service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /http-service .

# TEST
FROM build as test

# PRODUCTION
FROM alpine:latest as production

RUN apk --no-cache add ca-certificates
COPY --from=build /http-service ./
RUN chmod +x ./http-service

ENTRYPOINT ["./http-service"]

EXPOSE 8080