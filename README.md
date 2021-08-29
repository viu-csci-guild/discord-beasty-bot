# Beasty
Discord bot written in GoLang which accepts simple commands. It can be used for student role assignments on request.

[![Build](https://github.com/viu-csci-guild/discord-beasty-bot/actions/workflows/build.yaml/badge.svg)](https://github.com/viu-csci-guild/discord-beasty-bot/actions/workflows/build.yaml)
[![Linting](https://github.com/viu-csci-guild/discord-beasty-bot/actions/workflows/lint.yaml/badge.svg)](https://github.com/viu-csci-guild/discord-beasty-bot/actions/workflows/lint.yaml)
[![Maintainability](https://api.codeclimate.com/v1/badges/242bb98be3cb26be71d2/maintainability)](https://codeclimate.com/github/viu-csci-guild/discord-beasty-bot/maintainability)


# Environment Variables
- `TOKEN` contains discord token provided as part of bot creation
- `STUDENT_ROLE_ID` contains the discord ID # which represents the student role
- `START_ROOM_ID` determines where wakeup and shutdown messages for the bot are posted to
- `SERVER_ID` which server the bot connects and watches for events at startup

# Development
- Build and create the local dev container: `docker-compose up -d app`
- Attach to the container using VScode's extension: [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- Develop!
- you can run your binaries in the container: `cd cmd/beasty/ && go build && ./beasty`
