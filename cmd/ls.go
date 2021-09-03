package cmd

import (
	"fmt"
	"github.com/thomaszub/envls/internal"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all environmental variables",
	Long:  `Lists all environmental variables.`,
	Run:   lsMain,
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func lsMain(cmd *cobra.Command, args []string) {
	envReader := internal.NewDefaultReader()
	filterChain := internal.NewEmptyFilterHandler()
	filterChain.AppendFilter(internal.NewPrefixFilter("_"))
	formatter := internal.NewDelimiterFormatter(" -> ")

	envs := envReader.Read()
	filteredEnvs := filterChain.Filter(envs)
	formattedEnvs := formatter.Format(filteredEnvs)
	for _, formattedEnv := range formattedEnvs {
		fmt.Println(formattedEnv)
	}
}
