package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gitlab.tiendanimal.com/iskaypet/orders-consumer/model"
	"log"
	"time"
)

var ctx = context.Background()

func main() {
	var timeout int
	flag.IntVar(&timeout, "timeout", 500, "milliseconds")
	flag.Parse()
	log.Printf("... timeout simulation: %d", timeout)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     true,
	})

	// @TODO: @apineiro replace by AWS Elastic Cache
	keyValue := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "a300p011",
		DB:       0,
	})

	app.Post("/orders-consumer", func(c *fiber.Ctx) error {
		stringValue := string(c.Body())
		var messageDTO model.MessageDTO
		err := json.Unmarshal(c.Body(), &messageDTO)
		if err != nil {
			log.Fatal(err)
		}

		key := messageDTO.ID

		value, err := keyValue.Exists(ctx, key).Result()
		if err != nil {
			log.Println(err)
		}

		if value == 0 {
			err = keyValue.Set(ctx, key, stringValue, time.Minute*time.Duration(30)).Err()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Millisecond * time.Duration(timeout))
			log.Println(stringValue) // process the message
		} else {
			return c.SendString(stringValue)
		}

		return c.SendString(stringValue)
	})

	if err := app.Listen(":4000"); err != nil {
		log.Fatal(err)
	}
}
