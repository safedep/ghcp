package main

import (
	// This must always be the first import
	"github.com/safedep/ghcp/cmd/server"
	_ "github.com/safedep/ghcp/init"

	"fmt"

	"github.com/safedep/dry/log"
	"github.com/safedep/dry/obs"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:              "ghcp [OPTIONS] [COMMAND] [ARG...]",
		Short:            "GitHub Comments Proxy Service",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return fmt.Errorf("unknown command: %s", args[0])
		},
	}

	cobra.OnInitialize(func() {
		log.InitZapLogger(obs.AppServiceName("ghcp"), obs.AppServiceEnv("dev"))
	})

	cmd.AddCommand(server.NewServerCommand())

	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}
