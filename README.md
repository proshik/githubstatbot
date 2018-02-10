# GitHubStatBot

[![Build Status](https://travis-ci.org/proshik/githubstatbot.svg?branch=master)](https://travis-ci.org/proshik/githubstatbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/proshik/githubstatbot)](https://goreportcard.com/report/github.com/proshik/githubstatbot)
[![codecov](https://codecov.io/gh/proshik/githubstatbot/branch/master/graph/badge.svg)](https://codecov.io/gh/proshik/githubstatbot)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/proshik/githubstatbot/issues)

[Telegram bot](https://t.me/githubstatbot) which show GitHub statistic by languages, stars and forks. Written on GO.

## Run

1. Download [githubstatbot](https://github.com/proshik/githubstatbot/releases)
2. You need talk with [BotFather](https://telegram.me/botfather) and follow a few simple steps for register your bot and take access token(GITHUBSTATBOT_TELEGRAMTOKEN)
3. Go to [GitHub OAuth App](https://github.com/settings/developers) and create new OAuth Apps. You will have `Client ID`(GITHUBSTATBOT_GITHUBCLIENTID) and `Client Secret`(GITHUBSTATBOT_GITHUBCLIENTSECRET)
4. Export environment variables, taking on previous steps and start application:

```bash
# Required environment variables
$ export GITHUBSTATBOT_TELEGRAMTOKEN=
$ export GITHUBSTATBOT_GITHUBCLIENTID=
$ export GITHUBSTATBOT_GITHUBCLIENTSECRET=

# Not required environment variables
# export GITHUBSTATBOT_MODE=
# export GITHUBSTATBOT_PORT=
# export GITHUBSTATBOT_LOGDIR=
# export GITHUBSTATBOT_TLSDIR=
# export GITHUBSTATBOT_STATICFILESDIR=
# export GITHUBSTATBOT_DBPATH=
# export GITHUBSTATBOT_AUTHBASICUSERNAME=
# export GITHUBSTATBOT_AUTHBASICPASSWORD=

$ ./githubstatbot
```  

### Run with go

```bash
$ go build

$ export GITHUBSTATBOT_MODE=local 
$ export GITHUBSTATBOT_TELEGRAMTOKEN=telegram_token
$ export GITHUBSTATBOT_GITHUBCLIENTID=github_client_id
$ export GITHUBSTATBOT_GITHUBCLIENTSECRET=github_client_secret
$ ./githubstatbot
```

### Run with docker

NOT forget insert environment variable in command.

#### Simple start

Start bot only inline prompt command

```bash
docker build -t githubstatbot:latest .
docker run --rm -p 8080:8080 \
-e GITHUBSTATBOT_MODE='local' \
-e GITHUBSTATBOT_TELEGRAMTOKEN='' \
-e GITHUBSTATBOT_GITHUBCLIENTID='' \
-e GITHUBSTATBOT_GITHUBCLIENTSECRET='' \
--name githubstatbot githubstatbot:latest
```

#### Full start(production mode)

```bash
docker run --rm -p 8080:8080 -e GITHUBSTATBOT_TELEGRAMTOKEN='' \
-e GITHUBSTATBOT_GITHUBCLIENTID='' \
-e GITHUBSTATBOT_GITHUBCLIENTSECRET='' \
-e GITHUBSTATBOT_DBPATH='/app/data/database.db' \
--mount=type=bind,source="$(pwd)"/data,target=/app/data \
--name githubstatbot githubstatbot:latest
```

## Usage

See bot output information.

## TODO

- tests;
- internalization;
- change polling on webhook.
- build docker image on travis; 

## Patch 

Welcome!
