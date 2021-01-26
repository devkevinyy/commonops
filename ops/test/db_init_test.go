package test

import (
	"testing"

	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/models"
)

func TestInitDb(t *testing.T) {
	database.MysqlClient.AutoMigrate(&models.DmsInstance{})
}
