# gitoday

## Background
This is a straightforward CLI tool that fetches most excited GitHub repositories about today by programming language. It also offers AI-powered analysis of these repositories, summarizing the information in a reader-friendly format.
## Installation
1. source code
   - Clone the repository
   - Run `go mod tidy`
   - Run `go build -o gitoday`
   - Run `./gitoday -apiKey="xxx"`
## Usage
![Usage Example](https://github.com/winterfx/gitoday/blob/main/doc/usage.gif)
## Document
![](./doc/flow.png)

1. Choose the programming language you are interested in from the list in the CLI.
2. Fetch the trending repositories from GitHub.
3. Analyze the repositories using Azure openAI by dify service.
4. Display the analysis results in the CLI.

## Debugger
```bash
$ dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .
API server listening at: 127.0.0.1:43000
```
```bash
# Connect to it from another terminal
$ dlv connect 127.0.0.1:43000
```
## License
MIT
