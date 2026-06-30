package media

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/media/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/media/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cloudinary"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
)

type Module struct {
	handler *handler.MediaHandler
	auth    gin.HandlerFunc
}

func NewModule(
	cfg *config.Config,
	auth gin.HandlerFunc,
) (*Module, error) {

	client, err := cloudinary.New(cfg)

	if err != nil {
		return nil, err
	}

	service := service.New(client)

	h := handler.New(service)

	return &Module{

		handler: h,

		auth: auth,
	}, nil
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {

	media := api.Group("")
	media.Use(m.auth)

	media.POST("/uploads", m.handler.Upload)
}
