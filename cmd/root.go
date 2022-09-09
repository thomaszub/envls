package cmd

import (
	"fmt"

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
	var cfgs []internal.Config

	format, err := cmd.Flags().GetString(FORMATTER)
	if err != nil {
		return err
	}
	cfgs = append(cfgs, internal.WithFormat(format))

	listHidden, err := cmd.Flags().GetBool(ALL)
	if err != nil {
		return err
	}
	if listHidden {
		cfgs = append(cfgs, internal.WithAll())
	}

	searched, err := cmd.Flags().GetStringArray(SEARCH)
	if err != nil {
		return err
	}
	for _, s := range searched {
		cfgs = append(cfgs, internal.WithSearch(s))
	}

	pipeline, err := internal.NewPipeline(cfgs...)
	if err != nil {
		return err
	}
	res, err := pipeline.Execute()
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
