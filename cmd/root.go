package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags
	cfgFile  string
	password string
	rootCmd  = &cobra.Command{
		Use:   "gpctl",
		Short: "A tool for interacting with Greenplum",
		Long:  `A quick and "painless" way of interacting with Greenplum made by Bret Beatty`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to gpctl!")
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			prompt := &survey.Password{
				Message: "Please enter in your password",
			}
			survey.AskOne(prompt, &password, survey.WithValidator(survey.Required))
			viper.Set("password", password)

			hn := viper.Get("database.hostname")
			user := viper.Get("database.user")
			psswd := viper.Get("password")
			port := viper.GetString("database.port")
			dbname := viper.Get("database.database")

			connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", hn, port, user, psswd, dbname)

			viper.Set("connstring", connString)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gpctl)")

	rootCmd.AddCommand(sqlConn)
	rootCmd.AddCommand(qCmd)
	rootCmd.AddCommand(loadCmd)
}

func initConfig() {
	if cfgFile != "" {
		// use config file from flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home dir
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// search config in home dir with name .gpctl (no extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gpctl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
