package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ThreadBolt/threadbolt/pkg/generator"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ThreadBolt components",
	Aliases: []string{"gen", "g"},
}

var generateModelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Generate a new model",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := strings.Title(args[0])
		
		if err := generator.GenerateModel(modelName); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating model: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Generated model: %s\n", modelName)
	},
}

var generateControllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a new controller",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		controllerName := strings.Title(args[0])
		
		if err := generator.GenerateController(controllerName); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating controller: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Generated controller: %s\n", controllerName)
	},
}

func init() {
	generateCmd.AddCommand(generateModelCmd)
	generateCmd.AddCommand(generateControllerCmd)
}