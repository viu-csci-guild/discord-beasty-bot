# Beasty
Discord bot written in GoLang which accepts simple commands. It can be used for student role assignments on request.

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
