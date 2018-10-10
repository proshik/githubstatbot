# GitHubStatBot

[![Travis](https://travis-ci.org/proshik/githubstatbot.svg?branch=master)](https://travis-ci.org/proshik/githubstatbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/proshik/githubstatbot)](https://goreportcard.com/report/github.com/proshik/githubstatbot)
[![codecov](https://codecov.io/gh/proshik/githubstatbot/branch/master/graph/badge.svg)](https://codecov.io/gh/proshik/githubstatbot)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/proshik/githubstatbot/issues)

[Telegram bot](https://t.me/githubstatbot) which show GitHub statistic by languages, stars and forks. Written on GO.

## Run

1. Download [githubstatbot](https://github.com/proshik/githubstatbot/releases)
2. You need to talk with [BotFather](https://telegram.me/botfather) and follow the few simple steps for the register Telegram bot. So you will get access token(GITHUBSTATBOT_TELEGRAMTOKEN)
3. Go to [GitHub OAuth App](https://github.com/settings/developers) and create the new OAuth Apps. You will get `Client ID`(GITHUBSTATBOT_GITHUBCLIENTID) and `Client Secret`(GITHUBSTATBOT_GITHUBCLIENTSECRET)
4. Export the environment variables, from the previous steps and start the application:

```
# Required environment variables

# port, like: 8080
$ export PORT=

# database URL, like: postgres://postgres:password@localhost:5432/githubstatbot?sslmode=disable
$ export DATABASE_URL=

# telegram token, like: 47174:342lt;j34;lkgj;l3kgj
$ export GITHUBSTATBOT_TELEGRAMTOKEN=

# see to github account settings, like: b2fb4db59dj20g8d92d2
$ export GITHUBSTATBOT_GITHUBCLIENTID=

# see to the github account settings, like: 96d687e72c049sku1lf5567ca810cd09eaacbe6
$ export GITHUBSTATBOT_GITHUBCLIENTSECRET=

# Not required environment variables

# static files directory, default: ./static
$ export GITHUBSTATBOT_STATICFILESDIR=

# basic auth username, default: username
$ export GITHUBSTATBOT_AUTHBASICUSERNAME=

# basic auth password, default: password
$ export GITHUBSTATBOT_AUTHBASICPASSWORD=

$ ./githubstatbot
```  

### Run with go

```bash
$ go build

$ export PORT=port
$ export DATABASE_URL=database_url
$ export GITHUBSTATBOT_TELEGRAMTOKEN=telegram_token
$ export GITHUBSTATBOT_GITHUBCLIENTID=github_client_id
$ export GITHUBSTATBOT_GITHUBCLIENTSECRET=github_client_secret
$ ./githubstatbot
```

### Run with docker

You MUST NOT forget to insert the environment variables in command.

#### Simple start

Start bot only inline prompt command

```bash
$ docker build -t githubstatbot:latest .

$ docker run --rm -p 8080:8080 \
-e PORT='8080' \
-e DATABASE_URL='' \
-e GITHUBSTATBOT_TELEGRAMTOKEN='' \
-e GITHUBSTATBOT_GITHUBCLIENTID='' \
-e GITHUBSTATBOT_GITHUBCLIENTSECRET='' \
--name githubstatbot githubstatbot:latest
```

#### Full start(production mode)

```bash
$ docker run --rm -p 8080:8080 \
-e PORT='8080' \
-e DATABASE_URL='' \
-e GITHUBSTATBOT_TELEGRAMTOKEN='' \
-e GITHUBSTATBOT_GITHUBCLIENTID='' \
-e GITHUBSTATBOT_GITHUBCLIENTSECRET='' \
-e GITHUBSTATBOT_DBPATH='/app/data/database.db' \
--mount=type=bind,source="$(pwd)"/data,target=/app/data \
--name githubstatbot githubstatbot:latest
```

## Usage

See bot output information.

## TODO

- statistics by commits by week, quarter and year;
- support the notifications of user activities (statistics by commits) at the end of the week, quarter and year;
- increase the test coverage;
- add the internalization;
- change method of receiving messages from the Telegram servers, from polling to the webhook.

## Patch 

You are welcome!
