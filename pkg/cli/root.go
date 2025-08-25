package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "threadbolt",
	Short: "ThreadBolt is a Go web framework inspired by Spring Boot",
	Long: `ThreadBolt is a convention-over-configuration web framework for Go that provides
MVC architecture, built-in ORM, dependency injection, and CLI tools for rapid development.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(testCmd)
}