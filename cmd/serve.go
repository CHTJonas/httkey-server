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

var path string
var addr string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run web server",
	Long: "Runs the web server serving static files out of the directory in the given path, " +
		"or the current directory if none us given.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(path) == 0 {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Panicln(err)
			}
		}

		srv := server.NewWebserver(path, addr, 10*time.Second)
		srv.RegisterMiddleware(server.ServerHeaderMiddleware)
		srv.RegisterMiddleware(server.ProxyMiddleware)
		srv.RegisterMiddleware(server.DefaultLogMiddleware)

		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Println(err)
			}
		}()
		log.Println("Starting server")

		c := make(chan os.Signal, 1)
		// SIGQUIT or SIGTERM will not be caught.
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Println("Waiting up to 30 seconds for server to shutdown")
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Shutdown error:", err.Error())
		}
		log.Println("Goodbye!")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&path, "path", "p", "", "directory from which to serve files (default current directory)")
	serveCmd.Flags().StringVarP(&addr, "bind", "b", "localhost:8080", "address and port to bind to")
}
