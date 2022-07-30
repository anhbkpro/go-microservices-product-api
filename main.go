package main

import (
	"context"
	"github.com/anhbkpro/go-microservices-product-api/data"
	"github.com/anhbkpro/go-microservices-product-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	// create the handlers
	ph := handlers.NewProducts(l, v)

	// create a new serve mux and register the handlers
	r := mux.NewRouter()

	// handlers for API
	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.Update)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt)
	signal.Notify(signChan, os.Kill)

	sig := <-signChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}
