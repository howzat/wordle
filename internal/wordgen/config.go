package wordgen

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

const DictionaryBaseDirKey = "DICTIONARY_DIR"

type Config struct {
	BaseDir string `env:"DICTIONARY_DIR,required"`
}

func NewDictionaryConfig(ctx context.Context) (Config, error) {
	config := Config{}
	err := envconfig.Process(ctx, &config)
	return config, err
}
