package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mindlesstaucher/gini/api/v1/customer"
	"github.com/mindlesstaucher/gini/api/v1/material"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")

	publish()
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("/help/2", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

var benchHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	publish(client)
}

func sub(client mqtt.Client, serviceInstance int) {

	if serviceInstance > 0 {

		topic := fmt.Sprintf("/bench/%d", serviceInstance-1)
		token := client.Subscribe(topic, 1, benchHandler)
		token.Wait()
		fmt.Printf("Subscribed to topic %s", topic)
	}
}

func SetupMqtt(serviceInstance int) mqtt.Client {
	var broker = "10.0.4.74"
	var port = 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(fmt.Sprintf("pfreundt-service-%d", serviceInstance))
	opts.SetPassword("Pfreundt1979!")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client, serviceInstance)

	return client
}

func SetupDb(serviceInstance int) *gorm.DB {

	var db *gorm.DB
	var err error

	path := "db"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	path = fmt.Sprintf("db/%d", serviceInstance)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	path = fmt.Sprintf("db/%d/database.db", serviceInstance)

	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&material.MaterialModel{})
	db.AutoMigrate(&customer.CustomerModel{})

	return db
}

func RunWebApi(db *gorm.DB, serviceInstance int) {
	r := gin.Default()

	r.GET("/api/v1/customer", customer.GetCustomer(db))
	r.POST("/api/v1/customer", customer.PostCustomer(db))
	r.POST("/api/v1/customer/init", customer.InitCustomer(db))
	r.POST("/api/v1/customer/readBenchmark", customer.ReadBenchmark(db))
	r.POST("/api/v1/customer/updateBenchmark", customer.UpdateBenchmark(db))
	r.POST("/api/v1/customer/deleteBenchmark", customer.DeleteBenchmark(db))
	r.GET("/api/v1/material", material.MaterialGet(db))

	addr := fmt.Sprintf(":%d", 8080+serviceInstance)
	r.Run(addr)
}

func GetServiceInstance() int {
	var serviceInstance int = 0
	programName := os.Args[0]
	fmt.Println(programName)
	argLength := len(os.Args[1:])
	fmt.Printf("Arg length is %d\n", argLength)
	if argLength >= 1 {
		serviceInstance, _ = strconv.Atoi(os.Args[1])
	}
	fmt.Printf("serviceInstance is %d", serviceInstance)

	return serviceInstance
}

func main() {

	var db *gorm.DB

	serviceInstance := GetServiceInstance()

	SetupMqtt(serviceInstance)

	db = SetupDb(serviceInstance)

	RunWebApi(db, serviceInstance)

}
