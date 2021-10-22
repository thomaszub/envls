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

	lsCmd.Flags().BoolP("all", "a", false, "Show hidden environmental variables (starting with _) ")
	lsCmd.Flags().StringP("formatter", "f", "del,=", "Specifies the formatter (currently only del) with a comma separated list of configuration arguments")
}

func lsMain(cmd *cobra.Command, _ []string) error {
	envReader := internal.NewDefaultReader()
	filterChain := internal.NewEmptyFilterHandler()
	if !listHiddenVariables(cmd) {
		filterChain.AppendFilter(internal.NewNoPrefixFilter("_"))
	}
	formatter, err := dispatchFormatterFlag(cmd)
	if err != nil {
		return err
	}
	envs := envReader.Read()
	acceptedEnvs := filterChain.Accepted(envs)
	formattedEnvs := formatter.Format(acceptedEnvs)
	for _, formattedEnv := range formattedEnvs {
		fmt.Println(formattedEnv)
	}
	return nil
}

func listHiddenVariables(cmd *cobra.Command) bool {
	flag, err := cmd.Flags().GetBool("all")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error reading all flag:", err.Error())
		return false
	}
	return flag
}

func dispatchFormatterFlag(cmd *cobra.Command) (internal.Formatter, error) {
	flag, err := cmd.Flags().GetString("formatter")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error reading formatter flag:", err.Error())
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
