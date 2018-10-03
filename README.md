[![Build Status](https://travis-ci.org/fvdveen/mu2.svg?branch=master)](https://travis-ci.org/fvdveen/mu2) [![Go Report Card](https://goreportcard.com/badge/github.com/fvdveen/mu2)](https://goreportcard.com/report/github.com/fvdveen/mu2)[![codecov](https://codecov.io/gh/fvdveen/mu2/branch/master/graph/badge.svg)](https://codecov.io/gh/fvdveen/mu2)

# Mu2: another discord music bot

Mu2 is another discord music bot. It's defining feature is that it updates live to changes in the config.

## Index

* Prerequisites
* Install
* Configuration
* Running

## Prerequisites

Before installing make sure you have docker and docker-compose installed

## Install

Clone the git repository

```bash
git clone https://github.com/fvdveen/mu2
```

## Configuration

In the consul KV store in compose edit the bot/config key to have the following items: 
* your discord token
* your desired prefix
* your desired database and database settings
* your desired log settings

### An example config:

```json
{
  "log": {
  	"level": "info",
    "discord": {
      "level": "warn",
      "webhook": "MY_DISCORD_WEBHOOK"
    }
  },
  "database": {
  	"type": "postgres",
    "host": "postgres",
    "user": "mu2",
    "password": "mu2",
    "ssl": "disable"
  },
  "bot": {
  	"discord": {
      "token": "MY_DISCORD_TOKEN"
    },
    "prefix": "$"
  }
}
```

If you edit and save the config you should see the bot update itself within a couple of seconds.

## Running

Run docker-compose and you're good to go

```bash
docker-compose up
```