package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"inmemcache/pkg/cache"
)

type apiServer struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	cache cache.Cache
}

func main() {
	serverAddr := flag.String("serverAddr", "0.0.0.0", "HTTP server network address")
	serverPort := flag.Int("serverPort", 4000, "HTTP server network port")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	api := apiServer{
		infoLog:  infoLog,
		errorLog: errLog,
		cache:    cache.NewCache(),
	}

	// initialize http server
	srvAddr := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	srv := &http.Server{
		Addr:           srvAddr,
		Handler:        api.routes(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start http listener
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errLog.Fatal(err)
		}
	}()
	infoLog.Printf("http listener started at: %s", srvAddr)

	<-interrupt
	infoLog.Printf("initiating graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		api.cache.Flush()
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		infoLog.Fatalf("server shutdown error: %+v", err)
	}
}
