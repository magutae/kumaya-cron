package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/magutae/kumaya-cron/utils"
	tele "gopkg.in/telebot.v3"
)

type IssueKeyword []struct {
	Rank    int    `json:"rank"`
	Keyword string `json:"keyword"`
	Traffic int    `json:"traffic"`
}

type NewKeyword []struct {
	Keyword      string `json:"keyword"`
	SearchVolume int    `json:"searchVolume"`
}

func sendTelegram(text string) {
	config, err := utils.LoadConfig()
	utils.HandleErr(err)

	pref := tele.Settings{
		Token: config.TelegramToken,
	}

	b, err := tele.NewBot(pref)
	utils.HandleErr(err)

	c, err := b.ChatByID(int64(config.TelegramChatID))
	utils.HandleErr(err)

	b.Send(c, text, tele.ModeMarkdown, tele.NoPreview)

}

func sendIssueKeywords() {
	req, err := http.NewRequest("GET", "https://blackkiwi.net/api/service/keyword/issue-keywords", nil)
	utils.HandleErr(err)
	req.Header.Add("Referer", "https://blackkiwi.net/service/trend")

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.HandleErr(err)
	defer resp.Body.Close()

	var issueKeyword IssueKeyword
	json.NewDecoder(resp.Body).Decode(&issueKeyword)

	var sendMessage = "*블랙키위 인키 키워드*\n\n"
	for i := 0; i < 20; i++ {
		sendMessage += fmt.Sprintf("%d. %s\t", i, issueKeyword[i].Keyword)
		sendMessage += fmt.Sprintf("(%s)\t", humanize.Comma(int64(issueKeyword[i].Traffic)))
		sendMessage += fmt.Sprintf("[Link](https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=%s)\n", issueKeyword[i].Keyword)
	}

	sendTelegram(sendMessage)
}

func sendNewKeywords() {
	req, err := http.NewRequest("GET", "https://blackkiwi.net/api/service/keyword/new-keywords", nil)
	utils.HandleErr(err)
	req.Header.Add("Referer", "https://blackkiwi.net/service/trend")

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.HandleErr(err)
	defer resp.Body.Close()

	var data map[string]NewKeyword
	json.NewDecoder(resp.Body).Decode(&data)

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var sendMessage = "*블랙키위 신규 키워드*\n\n"
	for i := 0; i < len(keys); i++ {
		sendMessage += fmt.Sprintf("%s\n", keys[i])
		newKeyword := data[keys[i]]
		for j := 0; j < len(newKeyword); j++ {
			sendMessage += fmt.Sprintf(" - %s\t", newKeyword[j].Keyword)
			sendMessage += fmt.Sprintf("(%s)\t", humanize.Comma(int64(newKeyword[j].SearchVolume)))
			sendMessage += fmt.Sprintf("[Link](https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=%s)\n", newKeyword[j].Keyword)
		}
	}

	sendTelegram(sendMessage)
}

func main() {
	sendIssueKeywords()
	sendNewKeywords()
}
