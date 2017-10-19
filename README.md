# GitHubStatBot

[![Build Status](https://travis-ci.org/proshik/githubstatbot.svg?branch=master)](https://travis-ci.org/proshik/githubstatbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/proshik/githubstatbot)](https://goreportcard.com/report/github.com/proshik/githubstatbot)
[![codecov](https://codecov.io/gh/proshik/jalmew/branch/master/graph/badge.svg)](https://codecov.io/gh/proshik/githubstatbot)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/proshik/githubstatbot/issues)

Telegram bot written in GO. This bot show GitHub statistic by languages, stars and forks.

## Commands
```
[/auth]() - авторизация
[/language]() - статистика языков в репозиториях авторизованного пользователя
[/language]() *<repo_name>* - статистика языков заданного репозитория авторизованного пользователя
[/star]() - статистика по звездам в репозиториях авторизованного пользователя
[/star]() *<repo_name>* - статистика по звездам заданного репозитория авторизованного пользователя
[/fork]() - статистика по форкам пользовательских репозиториев авторизованного пользователя
[/fork]() *<repo_name>* - статистика по форкам заданного репозитория авторизованного пользователя
[/cancel]() - отмена авторизации
```
