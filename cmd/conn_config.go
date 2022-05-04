package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bab014/greenloader/functions"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sqlConn = &cobra.Command{
		Use:   "test-conn",
		Short: "Testing the connection to greenplum",
		Long:  `A way of checking to see if you have your connection properly setup for Greenplum`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Testing connection\n\n")

			connString := viper.GetString("connstring")
			conn, err := functions.Connect(connString)
			if err != nil {
				ec := color.New(color.FgHiRed).Add(color.Bold)
				ec.Fprintf(os.Stderr, "Error in connecting to Greenplum:\n%v\n", err)
				os.Exit(1)
			}
			defer conn.Close(context.Background())

			sc := color.New(color.FgHiGreen)
			sc.Println("Your conneciton to Greenplum is good to go.")
		},
	}
)
