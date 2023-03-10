package main

import (
	"encoding/json"
	"examples/caching"
	"flag"
	"log"
	"time"

	"examples/model"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var timeout int
	flag.IntVar(&timeout, "timeout", 500, "milliseconds")
	flag.Parse()
	log.Printf("... timeout simulation: %d", timeout)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     true,
	})

	appCache, err := caching.NewBuilder[string, model.MessageDTO]().
		Size(100).
		ExpireAfterWrite(time.Duration(5) * time.Minute).
		Build()

	app.Post("/orders-consumer", func(c *fiber.Ctx) error {
		stringValue := string(c.Body())
		var messageDTO model.MessageDTO
		err = json.Unmarshal(c.Body(), &messageDTO)
		if err != nil {
			log.Fatal(err)
		}

		key := messageDTO.ID

		value := appCache.GetIfPresent(key)

		if value.IsNone() {
			appCache.Put(key, messageDTO)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Millisecond * time.Duration(timeout))
			processMessage(stringValue, messageDTO)
		} else {
			return c.SendString(stringValue)
		}

		return c.SendString(stringValue)
	})

	if err = app.Listen(":4000"); err != nil {
		log.Fatal(err)
	}
}

func processMessage(stringValue string, messageDTO model.MessageDTO) {
	log.Println(stringValue) // process the message
	var orderDto model.OrderDTO
	err := json.Unmarshal([]byte(messageDTO.Msg), &orderDto)
	if err != nil {
		log.Println(err)
	}
	log.Printf("GET api.iskaypet.com/orders/%d -> { customer_id: 123, products: [001, 002] }", orderDto.ID)
	// log.Println("GET api.iskaypet.com/customers/123")
	// log.Println("GET api.iskaypet.com/products/001")
	// log.Println("GET api.iskaypet.com/products/002")
}
