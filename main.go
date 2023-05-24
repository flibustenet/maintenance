package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

// GCPLogStdout output json ready for gcp
// implement slog.Handler
type GCPLogStdout struct {
	attrs []slog.Attr
}

// Enabled implement slog.Handler for all levels
func (m *GCPLogStdout) Enabled(c context.Context, level slog.Level) bool {
	return true
}

// WithAttrs implement slog.Handler
func (m *GCPLogStdout) WithAttrs(attrs []slog.Attr) slog.Handler {
	nlog := &GCPLogStdout{}
	nlog.attrs = append(m.attrs, attrs...)
	return nlog
}

// WithGroup not implemented
func (m *GCPLogStdout) WithGroup(name string) slog.Handler {
	return m
}

type sGCP struct {
	Severity string            `json:"severity"`
	Message  string            `json:"message"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// Handle implement slog.Handler
func (m *GCPLogStdout) Handle(ctx context.Context, rec slog.Record) error {
	s := sGCP{}
	s.Severity = "INFO"
	s.Message = rec.Message
	switch rec.Level {
	case slog.LevelDebug:
		s.Severity = "DEBUG"
	case slog.LevelInfo:
		s.Severity = "INFO"
	case slog.LevelWarn:
		s.Severity = "WARNING"
	case slog.LevelError:
		s.Severity = "ERROR"
	}

	s.Labels = map[string]string{}
	for _, a := range m.attrs {
		s.Labels[a.Key] = a.Value.String()
	}
	rec.Attrs(func(a slog.Attr) {
		s.Labels[a.Key] = a.Value.String()
	})
	res, _ := json.Marshal(s)
	fmt.Printf("%s\n", res)
	return nil
}

func main() {
	slog.SetDefault(slog.New(&GCPLogStdout{}).With("maintenance", "maintenance"))
	/*
		opts := slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case "level":
					switch a.Value.String() {
					case "DEBUG":
						return slog.String("severity", "DEBUG")
					case "INFO":
						return slog.String("severity", "INFO")
					case "WARN":
						return slog.String("severity", "WARNING")
					case "ERROR":
						return slog.String("severity", "ERROR")
					}
				case "msg":
					a.Key = "message"
					return a
				}
				return a
			},
			Level: slog.LevelDebug,
		}
		slog.SetDefault(slog.New(opts.NewJSONHandler(os.Stdout)).With("maintenance", "maintenance"))

			opts := slog.HandlerOptions{
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == "level" {
						return slog.Attr{Key: "severity", Value: a.Value}
					}
					return a
				},
			}
			std := opts.NewJSONHandler(os.Stdout)
			slog.SetDefault(slog.New(std))
	*/
	slog.Info("message", "mylab", "mmm", "tylab", "ttt")
	slog.Warn("warning message", "mylab", "mmm", "tylab", "ttt")
	slog.Error("Error message", "err", fmt.Errorf("oups"))
	slog.Debug("Debug message")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", HelloServer)
	log.Println("Serve on :" + port)
	http.ListenAndServe(":"+port, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<!doctype html>
<title>Maintenance</title>
<style>
  body { text-align: center; padding: 150px; }
  h1 { font-size: 50px; }
  body { font: 20px Helvetica, sans-serif; color: #333; }
  article { display: block; text-align: left; width: 650px; margin: 0 auto; }
  a { color: #dc8100; text-decoration: none; }
  a:hover { color: #333; text-decoration: none; }
</style>

<article>
        <div title='%s'>Travaux en cours, merci de revenir un peu plus tard.</div>
</article>
`, os.Getenv("K_REVISION"))
}
