package main

import (
	"banking-app/internal/config"
	"banking-app/internal/model"
	"banking-app/internal/repository/banking"
	payment_system "banking-app/internal/service/payment-system"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	repo := banking.New(ctx, cfg.Banking.DSN)
	defer repo.Close()

	service := payment_system.New(ctx, repo)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      RegisterRouter(service),
	}

	go func() {
		log.Printf("service is running, avaliable on port: %d", 8080)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func RegisterRouter(service payment_system.Service) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid body request", http.StatusBadRequest)
			return
		}
		var transaction model.Transaction
		if err := json.Unmarshal(body, &transaction); err != nil {
			http.Error(w, "Invalid body request", http.StatusBadRequest)
			return
		}
		amount, err := transaction.Amount.Float64()
		if err != nil {
			http.Error(w, "Bad amount", http.StatusBadRequest)
			return
		}
		if err := service.Send(r.Context(), string(transaction.From), string(transaction.To), amount); err != nil {
			if errors.Is(err, banking.NotEnough) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, fmt.Sprintf("internal server error: %+v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		countStr := r.URL.Query().Get("count")
		count, err := strconv.Atoi(countStr)
		if err != nil {
			http.Error(w, "Invalid count parameter", http.StatusBadRequest)
			return
		}

		transactions, err := service.GetLast(r.Context(), count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			log.Printf("failed to encode transactions: %+v", err)
		}
	})

	router.HandleFunc("/api/wallet/{address}/balance", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address, ok := vars["address"]
		if !ok {
			http.Error(w, "bad address path", http.StatusBadRequest)
			return
		}
		balance, err := service.GetBalance(r.Context(), address)
		if err != nil {
			http.Error(w, "failed get balance", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(model.Balance{Balance: balance}); err != nil {
			log.Printf("failed to encode balance: %+v", err)
		}
	})
	return router
}
