package cmd

import (
	"errors"
	"fmt"
	"github.com/thomaszub/envls/internal"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all environmental variables",
	Long:  `Lists all environmental variables.`,
	RunE:  lsMain,
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
	lsCmd.Flags().StringP("formatter", "f", "del,=", "Specifies the formatter (currently only del) with a comma separated list of configuration arguments")
}

func lsMain(cmd *cobra.Command, _ []string) error {
	envReader := internal.NewDefaultReader()
	filterChain := internal.NewEmptyFilterHandler()
	filterChain.AppendFilter(internal.NewPrefixFilter("_"))
	formatter, err := dispatchFormatterFlag(cmd)
	if err != nil {
		return err
	}
	envs := envReader.Read()
	filteredEnvs := filterChain.Filter(envs)
	formattedEnvs := formatter.Format(filteredEnvs)
	for _, formattedEnv := range formattedEnvs {
		fmt.Println(formattedEnv)
	}
	return nil
}

func dispatchFormatterFlag(cmd *cobra.Command) (internal.Formatter, error) {
	flag, err := cmd.Flags().GetString("formatter")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading formatter flag:", err.Error())
		return nil, err
	}
	splitted := strings.Split(flag, ",")
	switch splitted[0] {
	case "del":
		if len(splitted) != 2 {
			return nil, errors.New(fmt.Sprintf("Formatter \"del\" takes exactly one configuration argument, e.g. del,=. Got: %s", flag))
		}
		return internal.NewDelimiterFormatter(splitted[1]), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown formatter: %s", flag))
	}

}
