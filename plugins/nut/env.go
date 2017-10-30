package nut

import "github.com/spf13/viper"

// Languages languages
func Languages() []string {
	return viper.GetStringSlice("languages")
}
