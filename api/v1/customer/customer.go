package customer

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"strconv"

	"time"
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
		offset := c.Query("offset")

		n, err := strconv.Atoi(limit)
		if err != nil {
			n = 1
		}

		o, err := strconv.Atoi(offset)
		if err != nil {
			o = 0
		}

		fmt.Printf("name:%v, limit:%v, offset:%v\n", name, limit, offset)

		query := fmt.Sprintf("%%%s%%", name)

		var customers []CustomerModel

		db.Limit(n).Offset(o).Where("name LIKE ?", query).Find(&customers)

		//fmt.Printf("%v", JSON(customers))

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

func createNCustomers(db *gorm.DB, n int64) {

	var i, c int64
	var batchSize int64 = 1000
	var batch []CustomerModel
	var cDto CustomerDto

	//Erstellen von 100000 Einträgen
	//1min18s für einzelne inserts
	//1.54 s für 1000er batches
	//2.48 s für 100er batches
	//11.57 s für 10er batches

	//Erstellen von 1 Mio Einträgen
	//13.13 s für 1000er batches

	if n < batchSize {
		batchSize = n
	}

	c = 0

	for i = 0; i < n; i++ {

		cDto = randomCustomer()
		batch = append(batch, CustomerModel{Code: cDto.Code, Name: cDto.Name, Price: cDto.Price})
		if int64(len(batch)) >= batchSize || i == n-1 {
			db.Create(&batch)
			c += 1
			fmt.Printf("Erstelle %vtes Batch mit %v Einträgen\n", c, batchSize)
			batch = batch[:0]
		}
	}

}

func InitCustomer(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var existing int64
		var required int64
		var requested int64

		start := time.Now()

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
				createNCustomers(db, required)
			}

			elapsed := time.Since(start)
			fmt.Printf("Creation took %s\n", elapsed)

			db.Model(&CustomerModel{}).Count(&existing)
			fmt.Printf("%v entries in db\n", existing)
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusNotAcceptable)
		}

	}
}
