package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thomaszub/envls/internal"
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
	rootCmd.Flags().BoolP("all", "a", false, "Show hidden environmental variables (starting with _) ")
	rootCmd.Flags().StringP("formatter", "f", "del,=", "Specifies the formatter with a comma separated list of configuration arguments. Possible values: del,DELIMITER and json,{compact,pretty}")
	rootCmd.Flags().StringArrayP("search", "s", []string{}, "Filter environmental variables by regex pattern matching names and values")
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
			return nil, fmt.Errorf("formatter \"del\" takes exactly one configuration argument, e.g. del,=. Got: %s", flag)
		}
		return internal.NewDelimiterFormatter(splitted[1]), nil
	case "json":
		if len(splitted) != 2 {
			return nil, jsonFormatterError(flag)
		}
		switch splitted[1] {
		case "compact":
			return internal.NewJsonFormatter(false), nil
		case "pretty":
			return internal.NewJsonFormatter(true), nil
		default:
			return nil, jsonFormatterError(flag)
		}
	default:
		return nil, fmt.Errorf("unknown formatter: %s", flag)
	}
}

func jsonFormatterError(flag string) error {
	return fmt.Errorf("formatter \"json\" takes exactly one configuration argument of compact or pretty, e.g. json,compact. Got: %s", flag)
}
