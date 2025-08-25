package cli

import (
	"fmt"
	"os"

	"github.com/ThreadBolt/threadbolt/pkg/generator"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [app-name]",
	Short: "Create a new ThreadBolt application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]

		if err := generator.CreateNewProject(appName); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully created ThreadBolt application: %s\n", appName)
		fmt.Printf("📁 Navigate to your project: cd %s\n", appName)
		fmt.Printf("🚀 Run your app: threadbolt run\n")
	},
}

func init() {
	newCmd.Flags().StringP("template", "t", "api", "Project template (api, web, minimal)")
}
