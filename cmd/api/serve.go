package api

import (
	"github.com/arfan21/backend-test/config"
	"github.com/arfan21/backend-test/internal/server"
	dbpostgres "github.com/arfan21/backend-test/pkg/db/postgres"
	dbredis "github.com/arfan21/backend-test/pkg/db/redis"
	"github.com/urfave/cli/v2"
)

func Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run the API server",
		Action: func(c *cli.Context) error {
			_, err := config.LoadConfig()
			if err != nil {
				return err
			}

			_, err = config.ParseConfig(config.GetViper())
			if err != nil {
				return err
			}

			dbRedis, err := dbredis.New()
			if err != nil {
				return err
			}

			dbPg, err := dbpostgres.NewPgx()
			if err != nil {
				return err
			}

			server := server.New(
				dbRedis,
				dbPg,
			)
			return server.Run()
		},
	}

}
