package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ThreadBolt/threadbolt/pkg/framework"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the ThreadBolt application server",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := framework.LoadApp()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading app: %v\n", err)
			os.Exit(1)
		}
		
		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			port = app.Config.GetString("server.port")
		}
		if port == "" {
			port = "8080"
		}
		
		fmt.Printf("🚀 Starting ThreadBolt application on port %s\n", port)
		
		if err := app.Start(port); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().StringP("port", "p", "", "Port to run the server on")
}