package cmd

import (
	"os"
	"os/exec"
	playlistcmd "yamp/cmd/playlist"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "yamp",
	Short:        "Yet another music player",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		tuiCommand := exec.Command(bunLocation(), tuiEntrypoint())
		tuiCommand.Stdin = os.Stdin
		tuiCommand.Stdout = os.Stdout
		tuiCommand.Stderr = os.Stderr
		tuiCommand.Env = append(os.Environ(), "INK_DISABLE_DEVTOOLS=1")

		return tuiCommand.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func tuiEntrypoint() string {
	entrypoint := os.Getenv("TUI_ENTRYPOINT")
	if entrypoint == "" {
		return "yamp-tui"
	} else {
		return entrypoint
	}
}

func bunLocation() string {
	bunLocation := os.Getenv("BUN_LOCATION")
	if bunLocation == "" {
		return "bun"
	} else {
		return bunLocation
	}
}

func init() {
	rootCmd.AddCommand(playlistcmd.Command())
}
