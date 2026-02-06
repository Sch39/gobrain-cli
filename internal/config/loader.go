package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load(dir string) (*Config, error) {
	if dir == "" {
		d, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = d
	}

	v := viper.New()
	v.SetConfigFile(filepath.Join(dir, "gob.yaml"))
	v.SetEnvPrefix("GOB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config

	hook := mapstructure.DecodeHookFunc(stringToSliceHook)

	if err := v.Unmarshal(&cfg, viper.DecodeHook(hook)); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

var stringToSliceHook = func(
	from reflect.Type,
	to reflect.Type,
	data interface{},
) (interface{}, error) {
	if from.Kind() != reflect.String || to != reflect.TypeOf([]string(nil)) {
		return data, nil
	}

	s, ok := data.(string)
	if !ok {
		return data, nil
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}, nil
	}

	return []string{s}, nil
}
