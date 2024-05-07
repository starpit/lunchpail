package subcommands

import (
	"lunchpail.io/pkg/cpu"

	"github.com/spf13/cobra"
)

func Newcmd() *cobra.Command {
	var namespaceFlag string
	var watchFlag bool

	var cmd = &cobra.Command{
		Use:   "cpu",
		Short: "Displays CPU utilization",
		Long:  "Displays CPU utilization",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cpu.UI(cpu.CpuOptions{Namespace: namespaceFlag, Watch: false})
		},
	}

	cmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Kubernetes namespace that houses your instance")
	cmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watches for changes in list of runs")

	return cmd
}

func init() {
	rootCmd.AddCommand(Newcmd())
}