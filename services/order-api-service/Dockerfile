FROM ubuntu:latest

ARG SERVICE_NAME="order-api-service"
WORKDIR /app/.bin/${SERVICE_NAME}

COPY ./.bin/${SERVICE_NAME}-linux /app/.bin/${SERVICE_NAME}/service
COPY ./.bin/${SERVICE_NAME}-migration-linux /app/.bin/${SERVICE_NAME}/migrate-db
COPY ./migrate/migrations/* /app/.bin/${SERVICE_NAME}/migrate/migrations/*
COPY ./entrypoint.sh /app/.bin/${SERVICE_NAME}/entrypoint.sh

ENV BIND_ADDRESS=":8080"

EXPOSE 8080

ENTRYPOINT [ "./entrypoint.sh"]
