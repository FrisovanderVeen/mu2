[![Build Status](https://travis-ci.org/fvdveen/mu2.svg?branch=master)](https://travis-ci.org/fvdveen/mu2) [![Go Report Card](https://goreportcard.com/badge/github.com/fvdveen/mu2)](https://goreportcard.com/report/github.com/fvdveen/mu2)[![codecov](https://codecov.io/gh/fvdveen/mu2/branch/master/graph/badge.svg)](https://codecov.io/gh/fvdveen/mu2)

# Mu2: another discord music bot

Mu2 is another discord music bot.

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

Build the executable

Create a config file

```bash
touch .env
```

## Configuration

Edit .env with the following items: 
* your discord token
* your desired prefix
* your desired log-level for stdout and discord
* your discord webhook

```env
MU2_DISCORD_TOKEN=YOUR_TOKEN
MU2_DISCORD_PREFIX=$
MU2_LOG_LEVEL_DISCORD=warn
MU2_LOG_LEVEL=info
MU2_LOG_WEBHOOK_DISCORD=YOUR_WEBHOOK
```

## Running

Run docker-compose and you're good to go

```bash
docker-compose up
```