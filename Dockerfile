# Lightweight alpine OS, weight only 5mb, everything else is Go environment
FROM golang AS builder
ARG HOST=${HOST}
ARG PORT=${PORT}
ARG DATABASE_URL=${DATABASE_URL}
ARG GITHUBSTATBOT_GITHUBCLIENTID=${GITHUBSTATBOT_GITHUBCLIENTID}
ARG GITHUBSTATBOT_GITHUBCLIENTSECRET=${GITHUBSTATBOT_GITHUBCLIENTSECRET}
ARG GITHUBSTATBOT_MODE=${GITHUBSTATBOT_MODE}
ARG GITHUBSTATBOT_TELEGRAMTOKEN=${GITHUBSTATBOT_TELEGRAMTOKEN}
ARG TZ=${TZ}
# Workdir is path in your docker image from where all your commands will be executed
WORKDIR /go/src/github.com/proshik/githubstatbot
# Add all from your project inside workdir of docker image
ADD . /go/src/github.com/proshik/githubstatbot
# Then run your script to install dependencies and build application
RUN CGO_ENABLED=0 go build -v
# Next start another building context
FROM alpine:3.6
# add certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# Change work directory
WORKDIR /app
# Copy static
COPY --from=builder /go/src/github.com/proshik/githubstatbot/static ./static
# Copy only build result from previous step to new lightweight image
COPY --from=builder /go/src/github.com/proshik/githubstatbot/githubstatbot .
# Expose port for access to your app outside of container
EXPOSE 80
# Starting bundled binary file
ENTRYPOINT [ "./githubstatbot" ]