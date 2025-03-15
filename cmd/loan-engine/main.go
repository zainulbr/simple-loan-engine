package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/libs/notification/mail"
	"github.com/zainulbr/simple-loan-engine/routes"
	"github.com/zainulbr/simple-loan-engine/settings"
)

func start(config *settings.Settings) error {
	ginModeMapper := map[string]string{
		settings.EnvDev:  gin.DebugMode,
		settings.EnvTest: gin.TestMode,
		settings.EnvProd: gin.ReleaseMode,
	}

	gin.SetMode(ginModeMapper[settings.Env()])

	router := gin.Default()
	base := router.Group(config.App.Server.APIBase)
	routes.RegisterLoanRoutes(base)
	routes.RegisterFileRoutes(base)
	server := &http.Server{
		Addr:    config.App.Server.HTTPAddress,
		Handler: router,
	}

	go func() {
		log.Println("server running on", config.App.Server.HTTPAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}

	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	server.Shutdown(context.Background())
	log.Println("Shutdown Server ...")
	return nil
}

func main() {
	config, err := settings.Load()
	if err != nil {
		log.Fatalln("error loading config", err)
	}
	fmt.Println(config)

	config.App.Version = fmt.Sprintf("v%d.%d", 0, 1)
	rootCmd := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: config.App.Name,
		Long:  config.App.Description,
		Run: func(cmd *cobra.Command, args []string) {
			// open connection
			err := pgsql.Open(config)
			if err != nil {
				log.Fatalln(err)
				return
			}
			err = mail.Open(config)
			if err != nil {
				log.Fatalln(err)
				return
			}

			// start server
			start(config)

			// close connection
			err = pgsql.Close()
			if err != nil {
				log.Fatalln(err)
				return
			}
			err = mail.Close()
			if err != nil {
				log.Fatalln(err)
				return
			}

		},
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "list all available configs",
		Long:  "List all available configurations in the application",
		Run: func(cmd *cobra.Command, args []string) {
			c, _ := json.Marshal(config)
			fmt.Println(string(c))
		},
	}
	rootCmd.AddCommand(configCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

}
