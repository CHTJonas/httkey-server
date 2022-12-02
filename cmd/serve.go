package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
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
		pwrBy := fmt.Sprintf("httkey/%s Go/%s (+https://github.com/CHTJonas/httkey-server)",
			version, strings.TrimPrefix(runtime.Version(), "go"))
		srv := server.NewWebserver(path, addr, 10*time.Second, pwrBy)
		srv.RegisterMiddleware(server.ServerHeaderMiddleware(pwrBy))
		srv.RegisterMiddleware(server.ProxyMiddleware)
		srv.RegisterMiddleware(server.LoggingMiddleware)

		log.Println("Starting server...")
		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalln("Startup error:", err.Error())
			}
		}()
		log.Println("Listening on", addr)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		signal.Notify(c, syscall.SIGQUIT)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		log.Println("Received shutdown signal!")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Println("Waiting for server to exit...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln("Shutdown error:", err.Error())
		}
		log.Println("Bye-bye!")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&addr, "bind", "b", "localhost:8080", "address and port to bind to")
}
