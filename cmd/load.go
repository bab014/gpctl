package cmd

import (
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bab014/greenloader/functions"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	truncate bool
	fle      string
	tbl      string
	loadCmd  = &cobra.Command{
		Use:   "load",
		Short: "A way of loading a CSV file into a table",
		Long:  "This command allows you to bulk load a CSV file in a Greenplum table",
		Run: func(cmd *cobra.Command, args []string) {

			// begin function
			// check that file is csv
			ext := strings.Split(fle, ".")[1]
			if ext != "csv" {
				ec := functions.ErrorColor()
				ec.Fprintf(os.Stderr, "You can only load csv files, you provided a %s\n", ext)
				os.Exit(1)
			}

			// Open csv file
			csvFile, err := os.Open(fle)
			if err != nil {
				ec := functions.ErrorColor()
				ec.Fprintf(os.Stderr, "The file %s does not exist\n", fle)
				os.Exit(1)
			}
			defer csvFile.Close()

			// Create database connection
			conn, _ := functions.Connect(viper.GetString("connstring"))

			// If `--truncate flag provided`
			if truncate {
				// Confirm with user that they actually wish to truncate
				doTruncate := false
				prompt := &survey.Confirm{
					Message: "Truncating. Are you sure? All data will be deleted",
				}
				survey.AskOne(prompt, &doTruncate)

				// User confirmed, run TRUNCATE
				if doTruncate {
					_, err := functions.Truncate(conn, tbl)
					if err != nil {
						ec := functions.ErrorColor()
						ec.Fprintf(os.Stderr, "An error occured when attempting to TRUNCATE %s: \n%v\n", tbl, err)
						os.Exit(1)
					}
					sc := color.New(color.FgHiGreen).Add(color.Bold)
					sc.Printf("Succesfully TRUNCATED %s.\n", tbl)
				} else {
					color.New(color.FgYellow).Print("Skipping TRUNCATION, will append instead\n\n")
				}
			}

			// Start `COPY` command
			color.New(color.FgHiYellow).Println("Beginning `COPY` command")
			res, err := functions.Importer(conn, csvFile, tbl)
			if err != nil {
				ec := functions.ErrorColor()
				ec.Fprintf(os.Stderr, "An error occured when Copying the csv file into %s:\n%v\n", tbl, err)
				os.Exit(1)
			}

			// Succesfull COPY Occurred
			color.New(color.FgHiGreen).Add(color.Bold).Printf("COPY succesfull, %d records loaded\n", res.RowsAffected())

		},
	}
)

func init() {
	loadCmd.Flags().BoolVar(&truncate, "truncate", false, "Truncate flag tells the load command whether or not to truncate the table before loading the file")
	loadCmd.Flags().StringVarP(&fle, "file", "f", "", "The file you wish to load")
	loadCmd.Flags().StringVarP(&tbl, "table", "t", "", "The table you wish to load data into. Formated (<schema>.<table>) eg. perf_cap.quotes_daily")
}
