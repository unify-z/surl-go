package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/config"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

func (ch *ConfigHandler) GetSiteConfig(c *gin.Context) {
	Success[map[string]any](c, map[string]any{
		"site_name":                 config.ConfigManagerInstance.Config.Site.SiteName,
		"allow_registration":        config.ConfigManagerInstance.Config.Site.AllowRegistration,
		"allow_guest_to_create_url": config.ConfigManagerInstance.Config.Site.AllowGuestToCreateURL,
	})
}
