package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/bab014/greenloader/functions"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	file   string
	output string
	query  string
	qCmd   = &cobra.Command{
		Use:   "query",
		Short: "Used for sending a query to Greenplum",
		Long:  `You can send a query to Greenplum and return the output to a file or stream to STDOUT`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Query Results\n\n")

			// Establish connection
			conn, err := functions.Connect(viper.GetString("connstring"))
			if err != nil {
				ec := color.New(color.FgHiRed).Add(color.Bold)
				ec.Fprintf(os.Stderr, "An error occurred in connecting to the database:\n%v\n", err)
			}
			defer conn.Close(context.Background())

			// Begin query logic
			if file != "" {
				// Read provided sql file
				sql, err := os.ReadFile(file)
				if err != nil {
					ec := functions.ErrorColor()
					ec.Fprintf(os.Stderr, "Could not open the sql file:\n%v\n", err)
					os.Exit(1)
				}

				// execute query
				rows, err := conn.Query(context.Background(), string(sql))
				if err != nil {
					ec := functions.ErrorColor()
					ec.Fprintf(os.Stderr, "An error occured when running the query:\n%v\n", err)
					os.Exit(1)
				}
				defer rows.Close()

				// Get column headers
				var columns []interface{}
				fieldDesc := rows.FieldDescriptions()
				for _, col := range fieldDesc {
					columns = append(columns, string(col.Name))
				}

				// Write resutls to csv
				if output != "" {
					ext := strings.Split(output, ".")
					if ext[1] != "csv" {
						ec := functions.ErrorColor()
						ec.Fprintf(os.Stderr, "Sorry, you can only output query results into a csv file, you tried a %s\n", ext[1])
						os.Exit(1)
					}

					f, err := os.Create(output)
					if err != nil {
						ec := functions.ErrorColor()
						ec.Fprintf(os.Stderr, "An error occured when opening %s: \n%v\v", output, err)
						os.Exit(1)
					}
					defer f.Close()

					w := csv.NewWriter(f)
					defer w.Flush()

					stringCols := make([]string, 0)
					for _, col := range columns {
						stringCols = append(stringCols, col.(string))
					}

					if err = w.Write(stringCols); err != nil {
						ec := functions.ErrorColor()
						ec.Fprintf(os.Stderr, "An error occured writing the column line to the output file: \n%v\n", err)
						os.Exit(1)
					}
					var rowNum int
					for rows.Next() {
						value, _ := rows.Values()
						stringValues, _ := functions.ConvertReturn(value)
						if err := w.Write(stringValues); err != nil {
							ec := functions.ErrorColor()
							ec.Fprintf(os.Stderr, "Error writing line %d to %s: \n%v\n", rowNum, output, err)
							os.Exit(1)
						}
						rowNum++
					}

					sc := color.New(color.FgHiGreen).Add(color.Bold)
					sc.Printf("Succesfully wrote the query results to %s\n", output)

					// Write results to STDIN
				} else {
					// Results loop
					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.AppendHeader(table.Row(columns))
					for rows.Next() {
						value, _ := rows.Values()
						strVals, _ := functions.ConvertReturn(value)
						infceSlice := functions.InterfaceSlice(strVals)
						t.AppendRow(table.Row(infceSlice))
						t.AppendSeparator()
					}
					t.Render()
				}
			}
		},
	}
)

func init() {
	qCmd.Flags().StringVarP(&file, "file", "f", "", "the sql file to be used for the query")
	qCmd.Flags().StringVarP(&output, "outputLocation", "o", "", "the location where you want the results file saved. If not provided, results will be sent to STDOUT")
	qCmd.Flags().StringVar(&query, "query", "", "the location where you want the results file saved. If not provided, results will be sent to STDOUT")
}
