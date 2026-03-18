package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rennanbadaro/cursor-usage-analyzer/internal/ui"
	"github.com/rennanbadaro/cursor-usage-analyzer/internal/usage"
	"github.com/spf13/cobra"
)

var version = "dev"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cursor-usage-analyzer",
		Short:   "Analyze Cursor usage CSV token consumption",
		Version: version,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			csvPath, err := ui.PromptCSVPath()
			if err != nil {
				if errors.Is(err, ui.ErrInputCancelled) {
					return ui.ErrInputCancelled
				}
				return err
			}

			if err := ui.ShowLoading("Processing CSV and building summary...", 1500*time.Millisecond); err != nil {
				return err
			}

			records, err := usage.ReadCSVFile(csvPath)
			if err != nil {
				return fmt.Errorf("error reading csv: %w", err)
			}

			summary := usage.Aggregate(records)
			report := ui.Render(summary)
			_, err = cmd.OutOrStdout().Write([]byte(report))
			if err != nil {
				return fmt.Errorf("write output: %w", err)
			}

			return nil
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	return cmd
}

func execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
