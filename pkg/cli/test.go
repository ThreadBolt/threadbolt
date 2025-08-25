package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧪 Running tests...")
		
		testCmd := exec.Command("go", "test", "./...")
		testCmd.Stdout = os.Stdout
		testCmd.Stderr = os.Stderr
		
		if err := testCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Tests failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ All tests passed")
	},
}