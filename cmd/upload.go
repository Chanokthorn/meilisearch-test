/*
Copyright © 2025 NAME HERE chanokthorn6@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"ms-tester/cmd/config"
	"ms-tester/meilisearch"
	"ms-tester/model"
	"ms-tester/runner"
	"ms-tester/storage/file_system"
	"time"

	"github.com/spf13/cobra"
)

type uploadMode string

const (
	uploadModeIterative uploadMode = "iterative"
	uploadModeBatch      uploadMode = "batch"
)

var (
	mode     uploadMode
	filePath string
	index    string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Read()

		fmt.Println("upload called on file:", filePath)

		loader := file_system.NewStreamLoader(file_system.WithReadBatchSize(100))
		loader.SetReadFile(filePath)
		loader.SetModel(model.Product{})

		ms := meilisearch.NewMeiliSearch(cfg.Host, cfg.MasterKey)

		if err := ms.CreateIndex(context.Background(), index, "id"); err != nil {
			fmt.Println(err)
			return
		}

		var wk runner.Worker
		switch mode {
		case uploadModeIterative:
			fmt.Println("upload mode iterative")
			wk = runner.NewIterativeWorker(ms).SetIndexUid(index)
		case uploadModeBatch:
			fmt.Println("upload mode bulk")
			wk = runner.NewBatchWorker(ms).SetBatchSize(10).SetIndexUid(index)
		default:
			fmt.Println("upload mode not found, exiting...")
			return
		}
		rn := runner.NewRunner().SetWorker(wk).SetWorkerAmount(3)

		start := time.Now()

		ctx := context.Background()

		latestTaskID, err := rn.Run(ctx, loader)
		if err != nil {
			err := fmt.Errorf("failed to run: %w", err)
			fmt.Println(err)
			return
		}

		if err := ms.WaitTaskDone(ctx, latestTaskID); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Upload completed")
		fmt.Printf("Time taken: %s\n", time.Since(start))

	},
}

func init() {
	uploadCmd.Flags().StringVarP((*string)(&mode), "mode", "m", "", "upload mode")
	uploadCmd.MarkFlagRequired("mode")
	uploadCmd.Flags().StringVarP(&filePath, "path", "p", "", "path directory")
	uploadCmd.MarkFlagRequired("mode")
	uploadCmd.Flags().StringVarP(&index, "index", "i", "", "index name")
	uploadCmd.MarkFlagRequired("index")
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
