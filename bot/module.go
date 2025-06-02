package bot

import (
	"github.com/luyanci/remake_bot/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	database.Module,

	fx.Provide(NewRemake),
	fx.Provide(NewMessageCounter),
	fx.Provide(NewHandler),
	fx.Invoke(func(handler *Handler) {
		handler.RegisterAll()
	}),
)
