# gitoday

## Background
This is a straightforward CLI tool that fetches the trending GitHub repositories for the current day by programming language. It also offers AI-powered analysis of these repositories, summarizing the information in a reader-friendly format.
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