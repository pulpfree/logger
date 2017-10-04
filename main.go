package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/pulpfree/auth"

	"github.com/goinggo/tracelog"
	"github.com/pulpfree/logger/model"
)

const (
	baseDir    = "."
	cfgDir     = baseDir + "/config"
	defaultEnv = "dev"
)

// Adapter type
type Adapter func(http.Handler) http.Handler

// Adapt method
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// helper method to set environment
func setEnv() (env string) {
	env = defaultEnv
	if e := os.Getenv("environment"); e != "" {
		env = e
	}
	return env
}

func logsHandler(e *auth.Env) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

			if r.Method == "OPTIONS" {
				return
			}

			if r.Method != "POST" {
				respondErr(w, http.StatusMethodNotAllowed, nil, "LoggerHandler")
				return
			}

			if r.Body == nil {
				err := errors.New("Missing body in request")
				respondErr(w, http.StatusBadRequest, err, "LoggerHandler")
				return
			}

			logger := new(model.Logger)
			if err := json.NewDecoder(r.Body).Decode(&logger.Body); err != nil {
				respondErr(w, http.StatusBadRequest, err, "LoginHandler")
				return
			}

			if err := logger.Record(e.DB); err != nil {
				respondErr(w, http.StatusBadRequest, err, "LoginHandler")
				return
			}

		})
	}
}

func main() {
	tracelog.Start(tracelog.LevelTrace)

	env := setEnv()
	e := &auth.Env{}
	e.Init(cfgDir, env)
	dbSession := e.SetDB()
	defer func() {
		tracelog.Info("Closing database connection...", "", "")
		dbSession.Close()
	}()

	r := http.NewServeMux()
	http.Handle("/", Adapt(r, logsHandler(e)))
	http.ListenAndServe(":3021", nil)

	tracelog.Stop()
}

func respondErr(w http.ResponseWriter, code int, err error, handlerNm string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	msg := fmt.Sprintf("%s", http.StatusText(code))
	if err != nil {
		msg += fmt.Sprintf(" - %s", err)
	}
	w.Write([]byte(msg))
	tracelog.Error(errors.New(http.StatusText(code)), handlerNm, msg)
}
