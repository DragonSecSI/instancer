package instancer

import (
	"time"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Instancer struct {
	Logger           zerolog.Logger
	HelmConfig       config.ConfigAppHelm
	KubernetesConfig config.ConfigAppKubernetes

	DB              *gorm.DB
	Prefix          string
	CleanupDuration time.Duration
}

type InstancerConfig struct {
	Name       string
	Repository string
	Version    string
	Values     []string
}
