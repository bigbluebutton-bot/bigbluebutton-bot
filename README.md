
# Big Blue Button Bot Prototype

## Overview
This project serves as a prototype for a bot designed specifically for Big Blue Button. The primary objective is to develop a highly performant bot that interfaces with the backend APIs of Big Blue Button without the need to initiate a headless browser.

The bot aims to utilize all the functionalities that a regular Big Blue Button user can access. If possible, the code should be written using Go (Golang).

## Prerequisites

### Node.js

#### Installation:

**Windows:**
1. Visit the [official Node.js download page](https://nodejs.org/en/download/).
2. Download the Windows Installer.
3. Run the installer and follow the on-screen instructions.

**Linux:**
1. Open your terminal.
2. Use [Node Version Manager](https://github.com/nvm-sh/nvm) to install Node.js:
   ```bash
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
   export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
   [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
   nvm install node
   ```

**macOS:**
1. Visit the [official Node.js download page](https://nodejs.org/en/download/).
2. Download the macOS Installer.
3. Run the installer and follow the on-screen instructions.

### Go (Golang)

#### Installation:

**Windows:**
1. Visit the [official Go download page](https://golang.org/dl/).
2. Download the Windows Installer.
3. Run the installer and follow the on-screen instructions.

**Linux:**
1. Open your terminal.
2. Use the following commands to download and install Go:
   ```bash
   wget https://dl.google.com/go/go1.xx.x.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.xx.x.linux-amd64.tar.gz
   ```

    Replace `1.xx.x` with the latest version number. You can find the latest version number on the [official Go download page](https://golang.org/dl/).

3. Add Go to your PATH:
   ```bash
   echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
   source ~/.bashrc
   ```

**macOS:**
1. Visit the [official Go download page](https://golang.org/dl/).
2. Download the macOS Installer.
3. Run the installer and follow the on-screen instructions.

### Configuration

Create a `config.json` file and insert the following configuration:
```json
{
    "bbb":{
       "api":{
          "url":"https://example.com/bigbluebutton/api/",
          "secret":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
          "sha":"SHA256"
       },
       "client":{
          "url":"https://example.com/html5client/",
          "ws":"wss://example.com/html5client/websocket"
       },
       "pad":{
          "url":"https://example.com/pad/",
          "ws":"wss://example.com/pad/"
       },
       "webrtc":{
          "ws":"wss://example.com/bbb-webrtc-sfu"
       }
    },
    "changeset": {
        "external": "false",
        "host": "0.0.0.0",
        "port": "5051"
    }
}
```
To retrieve the Big Blue Button secret/salt, execute the following command on the Big Blue Button server: `bbb-conf --secret`. For more details, refer to the [official documentation](https://docs.bigbluebutton.org/administration/bbb-conf/#--secret).

### Setup and Execution

1. Copy the `example.go` file from the `_example` directory.
2. Execute the following commands:
   ```go
   go mod init bbb-bot
   go mod tidy
   go run .
   ```

    Upon execution, the bot will create a new meeting room and join it. It will then initiate an English capture and write "Hello World" into it. Additionally, the bot will respond to messages sent in the main chat. For instance, if "ping" is written, the bot will reply with "pong".
