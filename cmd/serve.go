package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CHTJonas/httkey-server/pkg/server"
	"github.com/spf13/cobra"
)

var addr string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run web server",
	Long: "Runs the web server serving static files out of the directory in the given path, " +
		"or the current directory if none us given.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("httkey version", version)
		srv := server.NewWebserver(path, addr, 10*time.Second)
		srv.RegisterMiddleware(server.ServerHeaderMiddleware)
		srv.RegisterMiddleware(server.ProxyMiddleware)
		srv.RegisterMiddleware(server.LoggingMiddleware)

		log.Println("Starting server...")
		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Println("Startup error:", err.Error())
			}
		}()
		log.Println("Listening on", addr)

		c := make(chan os.Signal, 1)
		// SIGQUIT or SIGTERM will not be caught.
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Received shutdown signal!")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Println("Waiting for server to exit...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Shutdown error:", err.Error())
		}
		log.Println("Bye-bye!")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&addr, "bind", "b", "localhost:8080", "address and port to bind to")
}
