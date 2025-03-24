package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo

func main() {
	// Command line flags
	port := flag.String("p", "3000", "Server port")
	prefix := flag.String("prefix", "", "URL prefix for all requests")
	flag.Parse()

	// Get directory from args or use current directory
	dir := "./"
	if args := flag.Args(); len(args) > 0 {
		dir = args[0]
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		ch := make(chan struct{})
		go signalHandler(ch, shutdown)

		serve(*port, dir, *prefix)

		<-ch
		close(ch)
	}()

	wg.Wait()
}

func serve(port, dir, prefix string) {
	// Create Echo instance
	e = echo.New()
	e.HideBanner = true

	// Add middleware
	middleware.DefaultLoggerConfig.Format = "${time_rfc3339} ${remote_ip} | ${method} ${uri} ${status} ${latency_human}\n"
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Error: Directory '%s' does not exist", dir)
	}

	// Determine the URL path prefix
	urlPrefix := "/"
	if prefix != "" {
		// Ensure prefix starts with a slash
		if prefix[0] != '/' {
			prefix = "/" + prefix
		}
		// Remove trailing slash if present
		if len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
			prefix = prefix[:len(prefix)-1]
		}
		urlPrefix = prefix
	}

	// Serve static files with prefix
	e.Static(urlPrefix, dir)

	// If prefix is set, redirect root to prefix
	if urlPrefix != "/" {
		e.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusMovedPermanently, urlPrefix)
		})
	}

	// Start server
	if err := e.Start(":" + port); err != nil && !ignoreError(err) {
		log.Fatalf("failed to start the http server: %v", err)
	}
}

func shutdown() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	log.Println("Gracefully shutting down. Please wait...")
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown the http server: %v", err)
	}
}

func signalHandler(ch chan<- struct{}, fn func()) {
	c := make(chan os.Signal, 5)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Stop(c)
			fn()
			ch <- struct{}{}
			return
		}
	}
}

func ignoreError(err error) bool {
	return err == nil || err == io.EOF || err == syscall.EPIPE ||
		err == http.ErrServerClosed
}
