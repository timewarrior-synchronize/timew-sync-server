FROM golang:alpine AS build

ENV APP_HOME /go/src

# Install build dependencies
RUN apk add build-base

# Copy and build project
COPY . $APP_HOME

WORKDIR $APP_HOME
RUN go mod download
RUN go mod verify
RUN go build -o /bin/timew-server

# Assemble the resulting image
FROM alpine

RUN mkdir authorized_keys
COPY --from=build /bin/timew-server /bin/server

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
CMD [ "start" ]

