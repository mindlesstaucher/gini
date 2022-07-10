package material

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaterialModel struct {
	gorm.Model
	Code  string
	Name  string
	Price uint
}

func MaterialGet(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "material",
		})

	}
}
