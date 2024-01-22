FROM golang:1.21.0-alpine3.18 as build

WORKDIR /service

COPY go.mod .
COPY go.sum .
#
RUN go mod download

COPY cmd cmd
COPY internal  internal

RUN CGO_ENABLED=0 GOOS=linux go build -C cmd -o /app

FROM alpine:latest
COPY --from=build /app app

COPY db/migrations migrations

ARG PORT
ARG DB_USER
ARG DB_PASS
ARG DB_HOST
ARG DB_PORT
ARG DB_NAME
ARG PAGINATION_LIMIT
ARG DEBUG_MODE

ENV PORT=${PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASS=${DB_PASS}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_NAME=${DB_NAME}
ENV PAGINATION_LIMIT=${PAGINATION_LIMIT}
ENV DEBUG_MODE=${DEBUG_MODE}

CMD ["./app"]