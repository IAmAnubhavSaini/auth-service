package main

import (
	"auth-service/config"
	m "auth-service/middlewares"
	rh "auth-service/route-handlers"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.LoadJWTConfig()

	var rl = m.NewRateLimiter(config.RATE_LIMIT_GLOBAL, time.Second)

	http.HandleFunc("/register", m.LimitRate(rh.Register, rl))
	http.HandleFunc("/login", m.LimitRate(rh.Login, rl))
	http.HandleFunc("/protected", m.LimitRate(m.Auth(rh.ProtectedEndpoint), rl))

	// Create a channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Setup the main server
	server := &http.Server{Addr: config.PORT, Handler: nil}

	// Run the server in a goroutine so it doesn't block
	go func() {
		fmt.Println("Starting server at", config.PORT)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Listen: %s\n", err)
			fmt.Println("Trying secondary port:", config.PORT2)
			server.Addr = config.PORT2
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			} else {
				fmt.Println("Server started at", config.PORT2)
			}
		} else {
			fmt.Println("Server started at", config.PORT)
		}
	}()

	// Block until we receive a signal
	<-stop

	fmt.Println("Shutting down the server...")

	// Create a deadline to wait for the server to shut down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown Failed:%+v", err)
	} else {
		fmt.Println("Server exited properly")
	}
}
