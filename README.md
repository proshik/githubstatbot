# GitHubStatBot

[![Build Status](https://travis-ci.org/proshik/githubstatbot.svg?branch=master)](https://travis-ci.org/proshik/githubstatbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/proshik/githubstatbot)](https://goreportcard.com/report/github.com/proshik/githubstatbot)
[![codecov](https://codecov.io/gh/proshik/jalmew/branch/master/graph/badge.svg)](https://codecov.io/gh/proshik/githubstatbot)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/proshik/githubstatbot/issues)

Telegram bot which show GitHub statistic by languages, stars and forks. Written on GO.

## Demo

[Link](https://t.me/githubstatbot)

## Run

1. You need talk with [BotFather](https://telegram.me/botfather) and follow a few simple steps for register your bot.
2. Get access token for work with bot. You will have GITHUBSTATBOT_TELEGRAMTOKEN.
3. Go to [GitHub OAuth App](https://github.com/settings/developers) and registration you application. You will have GITHUBSTATBOT_GITHUBCLIENTID and GITHUBSTATBOT_GITHUBCLIENTSECRET.
4. Not you must export environment variables takes on previous steps:
 - GITHUBSTATBOT_TELEGRAMTOKEN
 - GITHUBSTATBOT_GITHUBCLIENTID
 - GITHUBSTATBOT_GITHUBCLIENTSECRET
5. Run githubstatbot.

### Run with go

```bash
$ go build
$ GITHUBSTATBOT_TELEGRAMTOKEN=<telegram_token> GITHUBSTATBOT_GITHUBCLIENTID=<github_client_id> GITHUBSTATBOT_GITHUBCLIENTSECRET=<github_client_secret> ./githubstatbot
```

### 5.2 Run with docker

NOT forget insert environment variable in command.

#### Linux, macOS

```bash
docker build -t githubstatbot:latest .
docker run --rm -p 8080:8080 -e GITHUBSTATBOT_TELEGRAMTOKEN='' \
-e GITHUBSTATBOT_GITHUBCLIENTID='' \
-e GITHUBSTATBOT_GITHUBCLIENTSECRET='' \
--name githubstatbot githubstatbot:latest
```

#### Windows

```bash
docker build -t githubstatbot:latest .
docker run --rm -p 8080:8080 -e GITHUBSTATBOT_TELEGRAMTOKEN='' ^
-e GITHUBSTATBOT_GITHUBCLIENTID='' ^
-e GITHUBSTATBOT_GITHUBCLIENTSECRET='' ^
--name githubstatbot githubstatbot:latest
```

## Usage

See bot output information.

## TODO

- tests;
- internalization;
- change polling method on webhook.
- build docker image on travis; 

## Patch 

Welcome!