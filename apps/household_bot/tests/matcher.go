package tests

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
)

func MsgConfigMatcher(want string) gomock.Matcher {
	return &messageConfigTextMatcher{tg.MessageConfig{Text: want}}
}

type messageConfigTextMatcher struct {
	got any
}

func (m *messageConfigTextMatcher) Matches(want any) bool {
	return want.(tg.MessageConfig).Text == m.got.(tg.MessageConfig).Text

}
func (m *messageConfigTextMatcher) String() string {
	return fmt.Sprintf("config text is: %s", m.got.(tg.MessageConfig).Text)
}
