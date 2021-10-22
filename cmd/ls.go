package cmd

import (
	"errors"
	"fmt"
	"github.com/thomaszub/envls/internal"
	"regexp"
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
	lsCmd.Flags().StringArrayP("search", "s", []string{}, "Filter environmental variables by regex pattern matching names and values")
}

func lsMain(cmd *cobra.Command, _ []string) error {
	envReader := internal.NewDefaultReader()
	filterChain := internal.NewEmptyFilterHandler()
	if err := applyAllFlag(cmd, &filterChain); err != nil {
		return err
	}
	if err := applySearchFlag(cmd, &filterChain); err != nil {
		return err
	}
	formatter, err := getFormatter(cmd)
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

func applyAllFlag(cmd *cobra.Command, filterHandler *internal.FilterHandler) error {
	listHiddenVars, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}
	if !listHiddenVars {
		filterHandler.AppendFilter(internal.NewNoPrefixFilter("_"))
	}
	return nil
}

func applySearchFlag(cmd *cobra.Command, filterHandler *internal.FilterHandler) error {
	searched, err := cmd.Flags().GetStringArray("search")
	if err != nil {
		return err
	}
	for _, s := range searched {
		regex, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		filterHandler.AppendFilter(internal.NewRegexFilter(regex))
	}
	return nil
}

func getFormatter(cmd *cobra.Command) (internal.Formatter, error) {
	flag, err := cmd.Flags().GetString("formatter")
	if err != nil {
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
