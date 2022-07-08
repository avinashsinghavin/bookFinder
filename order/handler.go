package order

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AggregationOrderHandler struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (oah AggregationOrderHandler) Init() {
	bookExtensionObject := BookOrderHttpExtension{
		Router: oah.Router,
		DB:     oah.DB,
	}
	bookExtensionObject.RegisterRoutes()
}
