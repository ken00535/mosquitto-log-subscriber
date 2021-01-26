package main

import (
	"fmt"
	"os"

	"mosquitto/log/pkg/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	hosts := []client.Host{}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	// level := viper.GetString("level")
	// output := viper.GetString("outputFile")
	err = viper.UnmarshalKey("hosts", &hosts)

	var rootCmd = &cobra.Command{
		Use:   "subscriber",
		Short: "Subscriber to subscribes mosquitto log",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(hosts)
			Subscribe(hosts[0])
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
