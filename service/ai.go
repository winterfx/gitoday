package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitoday/global"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	difyEndpoint = "https://api.dify.ai/v1/chat-messages"
)

var prompt = `
	你是一个GitHub代码分析师，请根据我给你的URL:%s分析出这个项目的信息。并按以下结构返回给我：
{
	"what":"",//这个项目是什么，请简要地尽量一句话概括它的功能。
	"why":["",""],//这个项目解决了哪些痛点,出于什么样的目的，请按数组的方式专业条理地列出。这一项突出解决了什么问题。
	"how":["",""],//这个项目是如何实现的，请列出它的使用的关键技术，请按数组的方式列出，这一项是突出技术。如果找不到细节，就请给出这种项目的通用设计。
	"other":["",""]//忽律这个项目的地址，找一些与这个项目相似的项目，知名度要高一些，列举几个他们的名字，按数组的方式列出。
}
example:
{
	"what":"Immich-Go is an open-source tool designed to streamline uploading large photo collections to your self-hosted Immich server. It is an alternative to the immich-CLI command that doesn't depend on NodeJS installation.",
	"why":["It solves the problem of handling massive archives downloaded from Google Photos using Google Takeout while preserving valuable metadata.",
	"It offers a simpler installation process than other tools, as it doesn't require NodeJS or Docker for installation.",
	"It discards any lower-resolution versions that might be included in Google Photos Takeout, ensuring the best possible copies on your Immich server."],
	"how":["Immich-Go uses the Immich API to interact with the Immich server.",
	"It supports uploading photos directly from your computer folders, folders tree and ZIP archives.",
	"It provides several options to manage photos, such as grouping related photos, controlling the creation of Google Photos albums in Immich, and specifying inclusion or exclusion of partner-taken photos."],
	"other":["rclone","gphotos-uploader-cli","gphotos-sync"]
],

	"how":[],
	"other":[]
}
当你写完之后，请再检查一下，确保你的回答是没有过多重复的内容和格式是否正确，请确保是json结构，请重新回答。
`

type ChatResponse struct {
	What  string   `json:"what"`
	Why   []string `json:"why"`
	How   []string `json:"how"`
	Other []string `json:"other"`
	Error error    `json:"error"`
}
type data struct {
	Event          string                 `json:"event"`
	ConversationId string                 `json:"conversation_id"`
	MessageId      string                 `json:"message_id"`
	CreatedAt      int64                  `json:"created_at"`
	TaskId         string                 `json:"task_id"`
	Id             string                 `json:"id"`
	Position       int                    `json:"position"`
	Thought        string                 `json:"thought"`
	Answer         string                 `json:"answer"`
	Observation    string                 `json:"observation"`
	Tool           string                 `json:"tool"`
	ToolLabels     map[string]interface{} `json:"tool_labels"`
	ToolInput      string                 `json:"tool_input"`
	MessageFiles   []interface{}          `json:"message_files"`
}

var apiKey string

func Init(key string) {
	apiKey = key
}

// Chat sends a request to the Dify API and prints the response
func Chat(ctx context.Context, repoUrl string) (*ChatResponse, error) {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*200)
		defer cancel()
	}
	if global.IsPreviewMode() {
		return &ChatResponse{
			What:  "Immich-Go is an open-source tool designed to streamline uploading large photo collections to your self-hosted Immich server. It is an alternative to the immich-CLI command that doesn't depend on NodeJS installation.",
			Why:   []string{"It solves the problem of handling massive archives downloaded from Google Photos using Google Takeout while preserving valuable metadata.", "It offers a simpler installation process than other tools, as it doesn't require NodeJS or Docker for installation.", "It discards any lower-resolution versions that might be included in Google Photos Takeout, ensuring the best possible copies on your Immich server."},
			How:   []string{"Immich-Go uses the Immich API to interact with the Immich server.", "It supports uploading photos directly from your computer folders, folders tree and ZIP archives.", "It provides several options to manage photos, such as grouping related photos, controlling the creation of Google Photos albums in Immich, and specifying inclusion or exclusion of partner-taken photos."},
			Other: []string{"rclone", "gphotos-uploader-cli", "gphotos-sync"},
		}, nil
	}
	requestBody, err := json.Marshal(map[string]interface{}{
		"inputs":          map[string]interface{}{},
		"query":           fmt.Sprintf(prompt, repoUrl),
		"response_mode":   "streaming",
		"conversation_id": "",
		"user":            "abc-123",
		"files": []map[string]string{
			{
				"type":            "image",
				"transfer_method": "remote_url",
				"url":             "https://cloud.dify.ai/logo/logo-site.png",
			},
		},
	})
	cr := &ChatResponse{}
	if err != nil {
		cr.Error = errors.Wrap(err, "json marshal error")
		return cr, err
	}

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "POST", difyEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		cr.Error = errors.Wrap(err, "create http request error")
		return cr, err
	}

	// Add headers
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		cr.Error = errors.Wrap(err, "http request error")
		return cr, err
	}
	defer resp.Body.Close()

	// Read the response body
	// Create a new buffered reader to handle the stream
	reader := bufio.NewReader(resp.Body)
	var answer string
	done := false
	for {
		if done {
			break
		}
		select {
		case <-ctx.Done():
			cr.Error = ctx.Err()
			return cr, err
		default:
			input, err := reader.ReadString('\n')
			if err != nil {
				// If the error is EOF, the stream ended normally
				if err == io.EOF {
					done = true
					break
				}
			}
			if strings.HasPrefix(input, "data: ") {
				input = strings.TrimPrefix(input, "data: ")
			}
			var d data
			err = json.Unmarshal([]byte(input), &d)
			if err != nil {
				continue
			}
			answer = answer + d.Answer
		}

	}
	err = json.Unmarshal([]byte(answer), cr)
	if err != nil {
		cr.Error = errors.Wrap(err, "json unmarshal error")
		return cr, err

	}
	return cr, nil
}
