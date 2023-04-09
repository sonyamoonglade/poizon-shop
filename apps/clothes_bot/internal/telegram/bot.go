package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	parseModeMarkdown = "markdown"
	parseModeHTML     = "html"
)

type Config struct {
	Token string
}

type bot struct {
	client *tg.BotAPI
}

func NewBot(config Config) (*bot, error) {
	client, err := tg.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	return &bot{
		client: client,
	}, nil
}

func (b *bot) GetUpdates() tg.UpdatesChannel {
	return b.client.GetUpdatesChan(tg.UpdateConfig{})
}

func (b *bot) Send(c tg.Chattable) (tg.Message, error) {
	return b.client.Send(c)
}

func (b *bot) CleanRequest(c tg.Chattable) error {
	_, err := b.client.Request(c)
	return err
}

func (b *bot) SendMediaGroup(c tg.MediaGroupConfig) ([]tg.Message, error) {
	return b.client.SendMediaGroup(c)
}

func (b *bot) Shutdown() {
	b.client.StopReceivingUpdates()
}
