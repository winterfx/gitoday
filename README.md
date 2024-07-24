# gitoday

## Background
This is a simple CLI tool that retrieves the current day's GitHub trending repositories by code language. It also supports AI analysis of these repositories.
## Install
1. source code
   - Clone the repository
   - Run `go build -o gitoday`
   - Run `./gitoday -apiKey="xxx"`
## Usage
## Document
## Debugger
```bash
$ dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .
API server listening at: 127.0.0.1:43000
```
```bash
# Connect to it from another terminal
$ dlv connect 127.0.0.1:43000
```