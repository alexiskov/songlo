package confreader

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func LoadConfig() (c ConfigEntity, err error) {
	if err = godotenv.Load(); err != nil {
		err = fmt.Errorf("env-file loading error: %w", err)
		return
	}
	if err = env.Parse(&c); err != nil {
		return
	}
	return
}
