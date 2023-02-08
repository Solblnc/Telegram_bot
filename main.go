package main

import (
	tgClient "Telegram_bot/clients/telegram"
	event_consumer "Telegram_bot/consumer/event-consumer"
	telegram "Telegram_bot/events/telegram"
	"Telegram_bot/storage/files"

	"flag"
	"log"
)

//6019123915:AAFov4mpfv0nQEaPtdbrGwreUi8pYhU05lk

const (
	thBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(tgClient.New(thBotHost, musttoken()), files.New(storagePath))

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("Service is stopped", err)
	}

}

func musttoken() string {
	token := flag.String("token-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is specified")
	}
	return *token
}
