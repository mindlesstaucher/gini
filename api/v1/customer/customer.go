package customer

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"strconv"
)

type CustomerDto struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type CustomerModel struct {
	gorm.Model
	Code  string
	Name  string
	Price float32
}

func randomName() (string, string) {

	firstNames := [...]string{"Ann-Katrin", "Almuth", "Giulia", "Felicitas", "Nicole", "Jule", "Klara", "Sara", "Laura", "Lina", "Alexandra", "Tabea", "Manuel", "Benedikt", "Mats", "Philipp", "Jérôme", "Bastian", "Mesut", "Toni", "Christoph", "Miroslav", "Thomas"}
	lastNames := [...]string{"Berger", "Schult", "Gwinn", "Rauch", "Anyomi", "Brand", "Bühl", "Däbritz", "Freigang", "Magull", "Popp", "Waßmuth", "Neuer", "Höwedes", "Hummels", "Lahm", "Boateng", "Schweinsteiger", "Özil", "Kroos", "Kramer", "Klose", "Müller"}

	i1 := rand.Intn(len(firstNames))
	i2 := rand.Intn(len(firstNames))
	i3 := rand.Intn(len(lastNames))

	n1 := firstNames[i1]
	n2 := firstNames[i2]
	n3 := lastNames[i3]

	return fmt.Sprintf("%s %s %s", n1, n2, n3), fmt.Sprintf("%s%s%s", n1[0:1], n2[0:1], n3[0:1])
}

func randomCustomer() CustomerDto {

	name, code := randomName()
	price := rand.Float32() * 100.0
	customer := CustomerDto{Name: name, Code: code, Price: price}

	return customer
}

func GetCustomer(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		name := c.Query("search")
		limit := c.Query("limit")

		n, err := strconv.Atoi(limit)
		if err != nil {
			n = 1
		}

		fmt.Printf("%v %v\n", name, limit)

		query := fmt.Sprintf("%%%s%%", name)

		var customers []CustomerModel

		db.Limit(n).Where("name LIKE ?", query).Find(&customers)

		c.JSON(http.StatusOK, customers)

	}
}

func PostCustomer(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		requestBody := CustomerDto{}
		err := c.Bind(&requestBody)

		if err != nil {
			panic(err)
		}

		db.Create(&CustomerModel{Code: requestBody.Code, Name: requestBody.Name, Price: requestBody.Price})

		c.Status(http.StatusOK)

	}
}

func InitCustomer(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var existing int64
		var required int64
		var requested int64
		var i int64

		db.Model(&CustomerModel{}).Count(&existing)

		n := c.Query("n")

		r, err := strconv.Atoi(n)
		requested = int64(r)

		required = requested - existing

		if required < 0 {
			db.Where("1 = 1").Delete(&CustomerModel{})
			required = requested
		}

		if err == nil {

			if required > 0 {
				for i = 0; i < required; i++ {
					rc := randomCustomer()
					db.Create(&CustomerModel{Code: rc.Code, Name: rc.Name, Price: rc.Price})
				}
			}

			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusNotAcceptable)
		}

	}
}
