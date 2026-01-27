package cms

import "github.com/adamkadda/arman/pkg/database"

type Config struct {
	Host  string `env:"HOST,required"`
	Port  string `env:"PORT,required"`
	Stage string `env:"STAGE" envDefault:"dev"`
	DB    *database.Config
}
