package pkgutil

import "github.com/arfan21/backend-test/config"

func GetPort() string {
	port := config.GetConfig().HttpPort
	if port != "" {
		return ":" + port
	}
	return ":8888"
}
