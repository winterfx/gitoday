package cmd

import (
	"flag"
	"fmt"
	"gitoday/global"
	"gitoday/service"
	"gitoday/ui/model"
	"io"
	"log"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

var helpText = fmt.Sprintf(`  
This is a command line tool that explore servicebus queue messages.  
  
The tool requires two command line arguments:
- env: This specifies the environment to be used. It can be either "int" or "stg".
- h: This prints the help text.

Example of usage:
./gitoday -mode=debug -preview=true
  
Please make sure that the queue name is valid according to Azure's naming rules.  
`)

func die() {
	flag.PrintDefaults()
	os.Exit(1)

}
func Execute() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s", helpText)
	}
	var mode string
	var apiKey string
	var preview bool

	flag.StringVar(&mode, "mode", "", "The environment to be used")
	flag.BoolVar(&preview, "preview", false, "Use fake data, not fetch from github")
	flag.Parse()
	apiKey = os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		die()
	}
	initLogger(mode)
	initGlobal(preview)
	initService(apiKey)
	slog.Info("Starting gitoday", slog.String("mode", mode), slog.String("apiKey", apiKey), slog.Bool("preview", preview))

	initModel()
}
func initGlobal(preview bool) {
	global.SetPreview(preview)
}
func initLogger(mode string) {
	if mode == "debug" {
		file, err := os.Create("./gitoday.log")
		if err != nil {
			panic(err)
		}
		if err := file.Truncate(0); err != nil {
			panic(err)
		}

		logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	} else {
		//disable slog
		logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
		slog.SetDefault(logger)
	}

}
func initService(apiKey string) {
	service.Init(apiKey)
}
func initModel() {
	p := tea.NewProgram(model.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Something went wrong %s", err)
	}
}
