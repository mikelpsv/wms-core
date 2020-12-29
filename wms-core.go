package wms_core

import (
	"fmt"
	"github.com/mikelpsv/wms-core/models"
)

func Version() {
	fmt.Println("Version 1.0.0")
}

func GetStorage() *models.Storage {
	storage := new(models.Storage)
	return storage
}
