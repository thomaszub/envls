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
	listHiddenVars, err := dispatchAllFlag(cmd)
	if err != nil {
		return err
	}
	if !listHiddenVars {
		filterChain.AppendFilter(internal.NewNoPrefixFilter("_"))
	}
	regexFilters, err := dispatchSearchFlag(cmd)
	if err != nil {
		return err
	}
	for _, filter := range regexFilters {
		filterChain.AppendFilter(filter)
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

func dispatchAllFlag(cmd *cobra.Command) (bool, error) {
	flag, err := cmd.Flags().GetBool("all")
	if err != nil {
		return flag, err
	}
	return flag, nil
}

func dispatchFormatterFlag(cmd *cobra.Command) (internal.Formatter, error) {
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

func dispatchSearchFlag(cmd *cobra.Command) ([]internal.Filter, error) {
	var filters []internal.Filter
	searched, err := cmd.Flags().GetStringArray("search")
	if err != nil {
		return nil, err
	}
	for _, s := range searched {
		regex, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		filters = append(filters, internal.NewRegexFilter(regex))
	}
	return filters, nil
}
