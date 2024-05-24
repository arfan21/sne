package main

import (
	"os"

	"github.com/arfan21/backend-test/cmd/api"
	migration "github.com/arfan21/backend-test/cmd/migrate"
	"github.com/urfave/cli/v2"
)

// @title sne API
// @version 1.0
// @description This is a sample server cell for sne Test API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	appCli := cli.NewApp()
	appCli.Name = "sne Test"
	appCli.Usage = "sne Test API"
	appCli.Commands = []*cli.Command{
		api.Serve(),
		migration.Root(),
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}
