package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/flohansen/sentinel/internal/file"
	"github.com/flohansen/sentinel/internal/proxy"
)

type App struct {
	version string
}

func NewApp(version string) *App {
	return &App{version}
}

func (app *App) PrintVersion(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("sentinel v%s\n", app.version)))
}

func (app *App) PrintHelp(w io.Writer) {
	app.PrintVersion(w)
	for _, line := range []string{
		"",
		"Usage: sentinel <command>",
		"",
		"Commands:",
		"  init     Create default configuration file",
		"  run      Run sentinel proxy",
		"  version  Print the binary version",
	} {
		w.Write([]byte(line + "\n"))
	}
}

func (app *App) Run(ctx context.Context, config Config) error {
	proxy := proxy.NewProxy(proxy.WithClient(&http.Client{}), proxy.WithConfig(config.Proxy))
	watcher := file.NewWatcher(proxy, config.Watch)

	errs := make(chan error)
	defer close(errs)

	go func() {
		if err := watcher.Start(ctx); err != nil {
			errs <- fmt.Errorf("watcher error: %s", err)
		}
	}()

	go func() {
		if err := proxy.Start(ctx); err != nil {
			errs <- fmt.Errorf("proxy error: %s", err)
		}
	}()

	return <-errs
}

func (a *App) Init() error {
	f, err := os.Create(".sentinel.json")
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(DefaultConfig); err != nil {
		return fmt.Errorf("error encoding config: %s", err)
	}

	return nil
}
