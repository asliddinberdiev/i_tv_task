package modules

import (
	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/user"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"gorm.io/gorm"
)

type Modules struct {
	User user.Service
}

func NewModules(cfg *config.Config, logger logger.Logger, db *gorm.DB) *Modules {
	userRepo := user.NewRepository(db)
	userSve := user.NewService(userRepo)

	return &Modules{
		User: userSve,
	}
}
