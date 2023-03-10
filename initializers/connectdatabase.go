package initializers

import (
	"os"

	"github.com/yurdigrfnn/api-todo-auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dns := os.Getenv("DATABASE")
	database, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(models.Todo{})
	database.AutoMigrate(models.User{})

	DB = database

}
