FROM golang:1.26-alpine3.24 AS build

ENV APP_HOME=/go/src

# Install build dependencies
RUN apk add build-base

# Copy and build project
COPY . $APP_HOME

WORKDIR $APP_HOME
RUN go mod download
RUN go mod verify
RUN go build -o /bin/timew-sync-server

# Assemble the resulting image
FROM alpine:3.24

RUN mkdir authorized_keys
COPY --from=build /bin/timew-sync-server /bin/timew-sync-server

EXPOSE 8080

ENTRYPOINT [ "/bin/timew-sync-server" ]
CMD [ "start" ]
