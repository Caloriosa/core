package tokenlib

import (
	"core/pkg/config"
	"strings"
)

func GetAppFromToken(token string) *string {
	parts := strings.Split(token, "/")
	if len(parts) != 2 {
		return nil
	}

	for _, app := range config.LoadedConfig.AppTokens {
		if app.Token == parts[1] && app.App == parts[0] {
			return &app.App
		}
	}
	return nil
}
