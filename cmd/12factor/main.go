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
	log := logrus.New()
	log.SetOutput(os.Stdout)

	log.Info("Starting the app...")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Infof("Receive request from %v", r.RemoteAddr)
	})

	serv := http.Server{
		Addr: net.JoinHostPort("", port),
		Handler: router,
	}

	go serv.ListenAndServe()

	log.Info("The app started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	log.Info("Stopping app..")

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := serv.Shutdown(timeout)
	if err != nil {
		log.Error("Error when shutdown app: %v", err)
	}

	log.Info("The app stopped")
}
