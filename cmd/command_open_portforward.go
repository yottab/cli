package cmd

import (
	"github.com/spf13/cobra"
)

// Open Command
var (
	openCmd = &cobra.Command{
		Use:   "open [Application.name]",
		Short: "Open an Application in the default browser",
		Long:  `This subcommand will open an application in the default browser in every os.`,
		Run:   appOpen}
)

// Portforward Command
var (
	portforwardCmd = &cobra.Command{
		Use:   "portforward [command] [name]",
		Short: "portforward to connect to an application running in a cluster",
		Long: `Portforward to connect to an application running in a cluster.
		This type of connection can be useful for database debugging`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	portforwardSrvCmd = &cobra.Command{
		Use:   "service [name]",
		Short: "portforward to connect to an application running in a cluster",
		Long: `Portforward to connect to an application running in a cluster.
		This type of connection can be useful for database debugging`,
		Run: srvPortforward}

	portforwardWorkerCmd = &cobra.Command{
		Use:   "worker [Application.name] [Worker.name]",
		Short: "port-forward to connect to an application worker running in a cluster",
		Long: `Port-forward to connect to an application worker running in a cluster.
		This type of connection can be useful for admin panels, monitoring tools`,
		Run: workerPortforward}
)

func init() {
	rootCmd.AddCommand(
		portforwardCmd,
		openCmd)

	portforwardCmd.AddCommand(
		portforwardSrvCmd,
		portforwardWorkerCmd)
}
