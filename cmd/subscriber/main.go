package main

import (
	"fmt"
	"io"
	"os"

	"mosquitto/log/pkg/client"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logrus.New()

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
	output := viper.GetString("outputFile")
	err = viper.UnmarshalKey("hosts", &hosts)

	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999-07:00",
	})

	var rootCmd = &cobra.Command{
		Use:   "subscriber",
		Short: "Subscriber to subscribes mosquitto log",
		Run: func(cmd *cobra.Command, args []string) {
			client.Subscribe(hosts[0])
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
