package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/vidmed/request"

	"github.com/pkg/errors"
)

var (
	configFileName = flag.String("config", "config.toml", "Config file name")

	requestService *request.Service
)

func init() {
	flag.Parse()
	config, err := NewConfig(*configFileName)
	if err != nil {
		request.GetLogger().Fatalf("ERROR loading config: %s\n", err.Error())
	}
	// InitLogger logging, logger goes first since other components may use it
	request.InitLogger(int(config.Main.LogLevel))

	// InitLogger Service
	requestService = request.NewRequestService()
}

func main() {
	numCPUs := runtime.NumCPU()
	request.GetLogger().Infof("CPUs count %d", numCPUs)
	runServer()
}

func runServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/request", recoverHandler(requestHandler))
	http.HandleFunc("/admin/requests", recoverHandler(viewsHandler))
	// todo perhaps remove global config
	hs := &http.Server{Addr: fmt.Sprintf("%s:%d", GetConfig().Main.ListenAddr, GetConfig().Main.ListenPort), Handler: nil}

	go func() {
		<-stop

		timeout := 15 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		fmt.Printf("Shutdown with timeout: %s\n", timeout)

		if err := hs.Shutdown(ctx); err != nil {
			request.GetLogger().Errorf("Error: %v\n", err)
		} else {
			request.GetLogger().Infof("Server stopped")
		}
		cancel()
	}()

	request.GetLogger().Infof("Listening on http://%s\n", hs.Addr)
	if err := hs.ListenAndServe(); err != http.ErrServerClosed {
		request.GetLogger().Error(err)
	}

	// stop service here
	requestService.Close()
}

func requestHandler(w http.ResponseWriter, _ *http.Request) {
	request.GetLogger().Debug("requestHandler")
	resp := requestService.GetRequest()

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write(resp.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewsHandler(w http.ResponseWriter, _ *http.Request) {
	request.GetLogger().Debug("viewsHandler")
	var resBytes []byte
	views := requestService.GetViews()
	if len(views) == 0 {
		resBytes = []byte("views are empty - please call /request")
	} else {
		resBytes = views.Bytes()
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write(resBytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func recoverHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if errR := recover(); errR != nil {
				err := errors.Errorf("panic: %+v", errR)
				request.GetLogger().Error(err)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		handler(w, r)
	}
}
