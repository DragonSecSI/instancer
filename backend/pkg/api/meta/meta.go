package meta

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers"
)

type MetaApi struct {
	DB     *gorm.DB
	Config config.ConfigApp
	Logger zerolog.Logger
}

func (rs MetaApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/conf", rs.GetConfig)

	return r
}

func (rs MetaApi) GetConfig(w http.ResponseWriter, r *http.Request) {
	helpers.Api.Response.Json(w, &rs.Logger, MetaConfigResponse{
		WebSuffix:    rs.Config.Meta.WebSuffix,
		SocketSuffix: rs.Config.Meta.SocketSuffix,
		SocketPort:   rs.Config.Meta.SocketPort,
	})
}
