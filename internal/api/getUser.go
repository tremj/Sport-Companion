package api

import (
	"fmt"
	"github.com/tremerj/Sport-Companion/database"
	"gorm.io/gorm"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var users []database.Users
	result := map[string]interface{}{}
	fmt.Println("MADE IT!")
	db.Find(&users).Take(&result)
	fmt.Println(result)
	for k, v := range result {
		fmt.Println(k, "+", v)
	}

}
