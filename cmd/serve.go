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

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run web server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.NewServer("/tmp", "127.0.0.1:8000", 10*time.Second)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateManifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateManifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
