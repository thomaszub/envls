package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/thomaszub/envls/internal/env"
	"github.com/thomaszub/envls/internal/filter"
	"github.com/thomaszub/envls/internal/format"
)

const (
	ALL       = "all"
	FORMATTER = "formatter"
	SEARCH    = "search"
)

var rootCmd = &cobra.Command{
	Use:   "envls",
	Short: "CLI tool for listing environmental variables",
	Long:  `CLI tool for listing environmental variables.`,
	RunE:  main,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolP(ALL, "a", false, "Show hidden environmental variables (starting with _) ")
	rootCmd.Flags().StringP(FORMATTER, "f", "del,=", "Specifies the formatter with a comma separated list of configuration arguments. Possible values: del,DELIMITER and json,{compact,pretty}")
	rootCmd.Flags().StringArrayP(SEARCH, "s", []string{}, "Filter environmental variables by regex pattern matching names and values")
}

func main(cmd *cobra.Command, _ []string) error {
	envReader := env.DefaultReader{}
	formatter, err := getFormatter(cmd)
	if err != nil {
		return err
	}
	envs := envReader.Read()
	filterHandler, err := makeFilterHandler(cmd)
	if err != nil {
		return err
	}
	acceptedEnvs := filterHandler.Accepted(envs)
	formattedEnvs, err := formatter.Format(acceptedEnvs)
	if err != nil {
		return err
	}
	fmt.Println(formattedEnvs)
	return nil
}

func getFormatter(cmd *cobra.Command) (format.Formatter, error) {
	flag, err := cmd.Flags().GetString(FORMATTER)
	if err != nil {
		return nil, err
	}
	return format.GetFormatter(flag)
}

func makeFilterHandler(cmd *cobra.Command) (*filter.FilterHandler, error) {
	var filters []filter.Filter
	all, err := allFlagFilters(cmd)
	if err != nil {
		return nil, err
	}
	filters = append(filters, all...)
	search, err := searchFlagFilters(cmd)
	if err != nil {
		return nil, err
	}
	filters = append(filters, search...)
	h := filter.FilterHandler{Filters: filters}
	return &h, nil
}

func allFlagFilters(cmd *cobra.Command) ([]filter.Filter, error) {
	listHiddenVars, err := cmd.Flags().GetBool(ALL)
	if err != nil {
		return nil, err
	}
	if listHiddenVars {
		return []filter.Filter{}, nil
	}
	f := filter.NoPrefixFilter{Prefix: "_"}
	return []filter.Filter{&f}, nil
}

func searchFlagFilters(cmd *cobra.Command) ([]filter.Filter, error) {
	searched, err := cmd.Flags().GetStringArray(SEARCH)
	if err != nil {
		return nil, err
	}
	var filters []filter.Filter
	for _, s := range searched {
		regex, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		f := filter.RegexFilter{Regex: regex}
		filters = append(filters, &f)
	}
	return filters, nil
}
