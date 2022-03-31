package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
)

type meta struct {
	RealName  string
	Name      string
	ID        string
	RealNames []string
}

func returnMeta(api *slack.Client, c slack.Channel) (metaInfo meta) {
	if c.IsMpIM {
		metaInfo.Name = c.ID + c.NameNormalized
	} else {
		user, err := api.GetUserInfo(c.Conversation.User)
		if err != nil {
			fmt.Println("error getting user: ", err)
		}
		metaInfo.Name = user.Name
		metaInfo.ID = user.ID
		metaInfo.RealName = user.RealName
	}
	return metaInfo
}

func main() {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	cursor := ""

	for {
		var userParams = slack.GetConversationsForUserParameters{
			UserID: os.Getenv("SLACK_USER"),
			Types:  []string{"im", "mpim"},
			Limit:  1000,
			Cursor: cursor,
		}

		channels, cursor, err := api.GetConversationsForUser(&userParams)

		if err != nil {
			fmt.Printf("%s\n", err)
		}

		for _, c := range channels {
			metaInfo := returnMeta(api, c)

			var page = ""

			err = os.Mkdir(metaInfo.Name, 0777)
			if err != nil {
				fmt.Println("error creating dir for user: ", metaInfo.Name)
			}

			conversationJson := make([]slack.Message, 0)
			limiter := time.Tick(1 * time.Second)

			for {
				<-limiter
				msgHistory, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
					ChannelID: c.Conversation.ID,
					Cursor:    page,
					Limit:     1000,
				})

				if err != nil {
					fmt.Println("error getting history: ", err)
				}

				for _, msg := range msgHistory.Messages {
					conversationJson = append(conversationJson, msg)
				}

				if msgHistory.ResponseMetaData.NextCursor == "" {
					type m struct {
						ConversationID string
						Name           string
						UserID         string
						Messages       []slack.Message
					}

					var d = m{
						ConversationID: c.Conversation.ID,
						Name:           metaInfo.RealName,
						UserID:         metaInfo.ID,
						Messages:       conversationJson,
					}

					js, err := json.MarshalIndent(d, "", "  ")

					if err != nil {
						fmt.Println("error writing json: ", err)
					}

					os.WriteFile(metaInfo.Name+"/messages.json", js, 0777)
					break
				}

				page = msgHistory.ResponseMetaData.NextCursor
			}
		}

		if cursor == "" {
			fmt.Println("no cursor present, breaking...")
			break
		}
	}
}
