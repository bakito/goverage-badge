package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bakito/goverage-badge/pkg/coverage"
	"github.com/bakito/goverage-badge/pkg/shield"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	fShieldConfig = "shield-config"
	fCoverageFile = "coverage-file"
	fLabel        = "label"
	fColor        = "color"
	fQuiet        = "quiet"
	vSeverity     = "severity"
	vTemplate     = "template"

	shellColor = "\033[1;33m%s\033[0m"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate badge json config",
	RunE: func(cmd *cobra.Command, args []string) error {
		severity := &shield.SeverityMap{}
		if err := viper.UnmarshalKey(vSeverity, severity); err != nil {
			return err
		}
		cf := viper.GetString(fCoverageFile)
		coverage, err := coverage.Calculate(cf)
		if err != nil {
			return err
		}

		badge := &shield.Badge{}
		if err := viper.UnmarshalKey(vTemplate, badge); err != nil {
			return err
		}
		badge.Setup(viper.GetString(fLabel), coverage, viper.GetString(fColor), severity)

		b, err := json.MarshalIndent(badge, "", "  ")
		if err != nil {
			return err
		}
		ioutil.WriteFile(viper.GetString(fShieldConfig), b, 0644)

		if !viper.GetBool(fQuiet) {
			cmd.Printf("Coverage is: %s\n", badge.Message)
			cmd.Printf("To add your badge to the readme as follows:\n")
			cmd.Println()
			cmd.Printf(shellColor, fmt.Sprintf("![Coverage](https://img.shields.io/endpoint?url=<url-to-your-%v>)\n", viper.GetString(fShieldConfig)))
			cmd.Println()
			cmd.Println("Visit: https://shields.io/endpoint for more details")
			cmd.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringP(fCoverageFile, "c", "coverage.out", "the coverage file to calculate the coverage value from (default: coverage.out)")
	_ = viper.BindPFlag(fCoverageFile, generateCmd.PersistentFlags().Lookup(fCoverageFile))

	generateCmd.PersistentFlags().StringP(fShieldConfig, "s", "coverage.json", "the shield.io json config file to be generated (default: coverage.json)")
	_ = viper.BindPFlag(fShieldConfig, generateCmd.PersistentFlags().Lookup(fShieldConfig))

	generateCmd.PersistentFlags().StringP(fLabel, "l", "coverage", "the badge label (default: coverage)")
	_ = viper.BindPFlag(fLabel, generateCmd.PersistentFlags().Lookup(fLabel))

	generateCmd.PersistentFlags().String(fColor, "", "the color for the badge. If not set, the color will be chosen based on severity.")
	_ = viper.BindPFlag(fColor, generateCmd.PersistentFlags().Lookup(fColor))

	generateCmd.PersistentFlags().BoolP(fQuiet, "q", false, "no output")
	_ = viper.BindPFlag(fQuiet, generateCmd.PersistentFlags().Lookup(fQuiet))
}
