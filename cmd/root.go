package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/thomaszub/envls/internal"
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
	formattedEnvs, err := formatter.Format(acceptedEnvs)
	if err != nil {
		return err
	}
	fmt.Println(formattedEnvs)
	return nil
}

func applyAllFlag(cmd *cobra.Command, filterHandler *internal.FilterHandler) error {
	listHiddenVars, err := cmd.Flags().GetBool(ALL)
	if err != nil {
		return err
	}
	if !listHiddenVars {
		f := internal.NewNoPrefixFilter("_")
		filterHandler.AppendFilter(&f)
	}
	return nil
}

func applySearchFlag(cmd *cobra.Command, filterHandler *internal.FilterHandler) error {
	searched, err := cmd.Flags().GetStringArray(SEARCH)
	if err != nil {
		return err
	}
	for _, s := range searched {
		regex, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		f := internal.NewRegexFilter(regex)
		filterHandler.AppendFilter(&f)
	}
	return nil
}

func getFormatter(cmd *cobra.Command) (internal.Formatter, error) {
	flag, err := cmd.Flags().GetString(FORMATTER)
	if err != nil {
		return nil, err
	}
	return internal.GetFormatter(flag)
}
