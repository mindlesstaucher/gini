package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mindlesstaucher/gini/api/v1/customer"
	"github.com/mindlesstaucher/gini/api/v1/material"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	var db *gorm.DB
	var err error

	db, err = gorm.Open(sqlite.Open("db/database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&material.MaterialModel{})
	db.AutoMigrate(&customer.CustomerModel{})

	r := gin.Default()

	r.GET("/api/v1/customer", customer.GetCustomer(db))
	r.POST("/api/v1/customer", customer.PostCustomer(db))
	r.POST("/api/v1/customer/init", customer.InitCustomer(db))

	r.GET("/api/v1/material", material.MaterialGet(db))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
