package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
	"strings"
)

type replaceString struct {
	Old      string
	New      string
	All      bool
	IsRegexp bool `mapstructure:"is_regexp"`

	compiledRegexp *regexp.Regexp
}

type config struct {
	ReplaceStringSlice []replaceString `mapstructure:"replace_strings"`
}

func main() {
	c, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range c.ReplaceStringSlice {
		if v.IsRegexp {
			compileRegexp(&c.ReplaceStringSlice[i])
		}
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		cs := cleanUpString(s.Text(), &c.ReplaceStringSlice)
		fmt.Printf("%v,%v\n", s.Text(), cs)
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}
}

func readConfig() (*config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	var c config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func compileRegexp(rs *replaceString) {
	rs.compiledRegexp = regexp.MustCompile(rs.Old)
}

func cleanUpString(s string, rs *[]replaceString) string {
	for _, v := range *rs {
		if v.IsRegexp {
			s = v.compiledRegexp.ReplaceAllString(s, v.New)
		} else if v.All {
			s = strings.ReplaceAll(s, v.Old, v.New)
		} else {
			s = strings.Replace(s, v.Old, v.New, 1)
		}
	}

	return s
}
