FROM ubuntu:latest

ARG SERVICE_NAME="currency-service"
WORKDIR /app/.bin/${SERVICE_NAME}

COPY ./.bin/${SERVICE_NAME}-linux /app/.bin/${SERVICE_NAME}/service

ENV BIND_ADDRESS=":8080"

EXPOSE 8080

ENTRYPOINT "./service"
