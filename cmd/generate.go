/*
Copyright © 2025 Chanokthorn Uerpairojkit chanokthorn6@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"ms-tester/generator"
	"ms-tester/model"
	"ms-tester/storage"

	"github.com/spf13/cobra"
)

var (
	count  int
	output string

	gn *generator.Generator
	st *storage.ProductStorage
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate filego",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called with", count, output)

		data, err := gn.Generate(count)
		if err != nil {
			fmt.Println(err)
			return
		}

		products := make([]model.Product, len(data))
		for i, d := range data {
			products[i] = *(d.(*model.Product))
		}

		err = st.SaveProduct(context.Background(), products, output)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	generateCmd.Flags().IntVarP(&count, "count", "c", 0, "number of files to generate")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "output directory")
	rootCmd.AddCommand(generateCmd)

	gn = generator.NewGenerator(&model.Product{})
	st = storage.NewProductStorage()

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
