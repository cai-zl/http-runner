/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行命令",
	Long:  `运行命令.`,
	Run: func(cmd *cobra.Command, args []string) {

		if runArgs.cron == "" && runArgs.second == -1 {
			fmt.Println("cron和time不能同时为空")
			os.Exit(2)
		}
		if runArgs.cron != "" && runArgs.second != -1 {
			fmt.Println("cron和time不能同时使用")
			os.Exit(2)
		}
		if runArgs.cron != "" {
			cronRunner()
		}
		if runArgs.second != -1 {
			timeRunner()
		}
	},
}

var runArgs struct {
	url    string
	cron   string
	second int64
}

func cronRunner() {
	c := cron.New(cron.WithLocation(time.UTC))
	_, err := c.AddFunc("@every 3s", func() {
		doRequest()
	})
	defer c.Stop()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	c.Start()
	select {}
}

func timeRunner() {
	duration := time.Second * time.Duration(runArgs.second)
	t := time.NewTimer(duration)
	defer t.Stop()
	go func() {
		for {
			<-t.C
			doRequest()
			t.Reset(duration)
		}
	}()
	select {}
}

func doRequest() {
	fmt.Printf("调用：%v\n", runArgs.url)
	client := resty.New()
	res, err := client.R().
		Get(runArgs.url)
	if err != nil {
		fmt.Printf("调用失败: %v\n", err.Error())
		os.Exit(2)
	}
	fmt.Printf("调用结果: %v\n", res.String())
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&runArgs.url, "url", "u", "", "http请求地址")
	_ = runCmd.MarkFlagRequired("url")
	runCmd.Flags().StringVarP(&runArgs.cron, "cron", "c", "", "cron表达式")
	runCmd.Flags().Int64VarP(&runArgs.second, "second", "s", -1, "time间隔")
}
