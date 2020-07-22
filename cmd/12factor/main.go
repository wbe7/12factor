package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)

func main()  {
	logrus.Info("Hello world!")

	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("Port is not set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		logrus.Infof("Receive request from %v", r.RemoteAddr)
	})

	serv := http.Server{
		Addr: net.JoinHostPort("", port),
		Handler: router,
	}

	go serv.ListenAndServe()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	serv.Shutdown(timeout)
}
