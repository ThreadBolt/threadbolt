package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ThreadBolt/threadbolt/pkg/framework"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := framework.LoadApp()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading app: %v\n", err)
			os.Exit(1)
		}
		
		if err := app.RunMigrations(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running migrations: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Migrations completed successfully")
	},
}