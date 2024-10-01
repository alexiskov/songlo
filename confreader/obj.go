package confreader

type (
	ConfigEntity struct {
		Server SrverEntity
		DMS    DatabaseManagementSystemEntity
	}

	SrverEntity struct {
		Port uint16 `env:"SERVER_PORT"`
	}

	DatabaseManagementSystemEntity struct {
		Host     string `env:"DB_HOST"`
		Username string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Port     uint16 `env:"DB_PORT"`
		DBname   string `env:"DB_NAME"`
	}
)
