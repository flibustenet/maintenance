package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"log/slog"

	"cloud.google.com/go/compute/metadata"
	"github.com/google/uuid"
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
	rec.Attrs(func(a slog.Attr) bool {
		s.Labels[a.Key] = a.Value.String()
		return true
	})
	res, _ := json.Marshal(s)
	fmt.Printf("%s\n", res)
	return nil
}

var projectID string

func main() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID, _ = metadata.ProjectID()
	}
	fmt.Println(Entry{
		Severity:  "NOTICE",
		Message:   "GOOGLE_CLOUD_PROJECT=" + projectID,
		Component: "arbitrary-property",
		//          Trace:     trace,
	})
	/*
		fmt.Println(Entry{
			Severity:  "INFO",
			Message:   "This is the default display field.",
			Component: "arbitrary-property",
			Trace:     "1",
		})
		fmt.Println(Entry{
			Severity:  "INFO",
			Message:   "deux",
			Component: "arbitrary-property",
			Trace:     "1",
		})
	*/

	//	slog.SetDefault(slog.New(&GCPLogStdout{}).With("maintenance", "maintenance"))
	opts := &slog.HandlerOptions{
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
			case "time":
				return slog.Attr{}
			case "msg":
				a.Key = "message"
				return a
			case "trace":
				return slog.String("logging.googleapis.com/trace",
					fmt.Sprintf("projects/%s/traces/%s", projectID, a.Value.String()))
			}
			return a
		},
		Level: slog.LevelDebug,
	}

	//trace := uuid.NewString()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, opts)))
	sj := slog.With("starting", "starting",
		"trace", uuid.NewString())
	sj.Info("message", "mylab", "mmm", "tylab", "ttt")
	sj.Warn("warning message", "mylab", "mmm", "tylab", "ttt")
	sj.Error("Error message", "err", fmt.Errorf("oups"))
	sj.Debug("Debug message")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", HelloServer)
	sj.Info("Serve on :" + port)
	http.ListenAndServe(":"+port, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	trace := ""
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = traceParts[0]
	}

	slog.Info("ici", "trace", trace)
	slog.Info("la", "trace", trace)
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

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}
