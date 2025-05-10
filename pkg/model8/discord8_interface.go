package model8

import "github.com/projectdiscovery/nuclei/v3/pkg/output"

type Model8Discord8Interface interface {
	InitialiseChannelID() error
	SetWebHook(string, string, string)
	SetBot(string)
	SetChatMessages([]output.ResultEvent) error
	GetChannelID() string
	GetBotToken() string
	GetChatMessages() []CustomMessagesHost
	AddChatMessages(CustomMessagesHost) []CustomMessagesHost
}
