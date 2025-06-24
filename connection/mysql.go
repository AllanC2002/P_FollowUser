package connection

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	user := os.Getenv("DBU_USER")
	password := os.Getenv("DBU_PASSWORD")
	host := os.Getenv("DBU_HOSTIP")
	port := os.Getenv("DBU_PORT")
	dbname := os.Getenv("DBU_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
