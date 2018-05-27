# Mu2: another discord music bot

Mu2 is another discord music bot.

## Index

* Prerequisites
* Install
* Configuration

## Prerequisites

Before installing make sure you have vgo installed

```bash
go get -u golang.org/x/vgo
```

## Install

Clone the git repository

```bash
git clone https://github.com/fvdveen/mu2
```

Build the executable

```bash
vgo build
```

Create a config file

```bash
touch config.yaml
```

## Configuration

Edit config.yaml with the following items: your discord token, your desired prefix and your discord bot's client id

```yaml
bot:
  token: "YOUR_DISCORD_TOKEN"
  prefix: "YOUR_PREFIX"
  invite-link: "https://discordapp.com/oauth2/authorize?client_id=YOUR_DISCORD_CLIENT_ID&scope=bot"
```

Thats it you can now run the bot