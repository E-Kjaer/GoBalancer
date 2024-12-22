package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordBot struct {
	Url string
}

type DiscordMessage struct {
	Username  string         `json:"username"`
	AvatarUrl string         `json:"avatar_url"`
	Content   string         `json:"content"`
	Embeds    []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Author      DiscordAuthor  `json:"author"`
	Title       string         `json:"title"`
	Url         string         `json:"url"`
	Description string         `json:"description"`
	Color       int32          `json:"color"`
	Fields      []DiscordField `json:"fields"`
}

type DiscordAuthor struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func (bot *DiscordBot) SendOfflineMessage(msg string) {
	bot.sendMessage(msg, 11734542)
}

func (bot *DiscordBot) SendOnlineMessage(msg string) {
	bot.sendMessage(msg, 1094416)
}

func (bot *DiscordBot) sendMessage(msg string, color int32) {
	embeds := []DiscordEmbed{}
	embed := DiscordEmbed{
		Author: DiscordAuthor{
			Name:    "",
			Url:     "",
			IconUrl: "",
		},
		Title:       msg,
		Url:         "",
		Description: "",
		Color:       color,
		Fields:      nil,
	}
	embeds = append(embeds, embed)
	message := DiscordMessage{
		Username:  "Load Balancer",
		AvatarUrl: "",
		Content:   "",
		Embeds:    embeds,
	}

	j, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = http.Post(bot.Url, "application/json", bytes.NewBuffer(j))
}
