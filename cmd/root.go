package cmd

import (
	"fmt"
	"github.com/raojinlin/apollo-client/apollo"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "apollo-client",
	Short: "apollo-client is a apollo config client",
	Long:  `A fast and easy to use command line tool for apollo`,
	Run: func(cmd *cobra.Command, args []string) {
		var server = viper.GetString("server")
		var cluster = viper.GetString("cluster")
		var appId = viper.GetString("appId")
		var cacheDir = viper.GetString("output")
		var namespaces = viper.GetStringSlice("namespaces")

		_, err := apollo.PullConfigAndSave(cacheDir, server, appId, cluster, namespaces)
		if err != nil {
			fmt.Println("Pull config go error: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("%v configuration saved.\n", namespaces)
		}

		if viper.GetBool("watch") {
			err = watch(server, appId, cluster, cacheDir, namespaces, &Notify{
				Script: viper.GetString("notify"),
				Url:    viper.GetString("notifyUrl"),
			})
			if err != nil {
				fmt.Println("watch change error, ", err.Error())
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cfgFile string
var server string
var appId string
var cluster string
var namespaces []string
var cacheDir string
var notify string
var notifyUrl string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.yaml", "Config file")
	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "Config server eg. http://192.168.71.111:8081/")
	rootCmd.PersistentFlags().StringVarP(&appId, "appId", "i", "", "App id")
	rootCmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", "default", "Cluster name")
	rootCmd.PersistentFlags().StringVarP(&cacheDir, "output", "o", "./", "Output path.")
	rootCmd.PersistentFlags().StringVarP(&notify, "notify", "n", "", "Notify command execute if changed.")
	rootCmd.PersistentFlags().StringVarP(&notifyUrl, "notifyUrl", "u", "", "Push server if changed.")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	rootCmd.PersistentFlags().Bool("watch", false, "Listen for configuration change")
	rootCmd.PersistentFlags().StringArray("namespace", namespaces, "App namespace, (default: [application])")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("appId", rootCmd.PersistentFlags().Lookup("appId"))
	viper.BindPFlag("cluster", rootCmd.PersistentFlags().Lookup("cluster"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("notify", rootCmd.PersistentFlags().Lookup("notify"))
	viper.BindPFlag("namespaces", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.BindPFlag("watch", rootCmd.PersistentFlags().Lookup("watch"))
	viper.BindPFlag("notifyUrl", rootCmd.PersistentFlags().Lookup("notifyUrl"))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
