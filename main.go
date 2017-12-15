package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func main() {
	if err := Main(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Main(args []string) error {
	token := args[1]
	channels, err := fetchChannels(token)
	if err != nil {
		return errors.Wrap(err, "failed to fetchChannels")
	}

	result := make(map[string]uint64, 10000)

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	since := time.Date(2017, time.January, 1, 0, 0, 0, 0, jst)

	for i, c := range channels {
		log.Printf("%d/%d %s", i, len(channels), c.Name)
		if strings.HasPrefix(c.Name, "noti_") {
			log.Println("ignore notification channel")
			continue
		}
		messages, err := fetchMessages(c.ID, since, token)
		if err != nil {
			return errors.Wrap(err, "failed to fetchMessages")
		}
		for _, m := range messages {
			for _, r := range m.Reactions {
				result[r.Name] += uint64(r.Count)
			}
		}
	}

	b, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}
	if err := ioutil.WriteFile(strconv.Itoa(int(time.Now().UnixNano()))+".json", b, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}

func fetchChannels(token string) (channels []Channel, err error) {
	channels = make([]Channel, 0, 1000)
	u, err := url.Parse("https://slack.com/api/channels.list")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	v := u.Query()
	v.Set("token", token)
	v.Set("exclude_members", "true")

	cursor := ""
	for {
		if cursor != "" {
			v.Set("cursor", cursor)
		}

		u.RawQuery = v.Encode()

		resp := ChannelListResponse{}
		if err = httpGet(u, &resp); err != nil {
			return nil, errors.Wrap(err, "failed to http get")
		}
		if !resp.Ok {
			return nil, errors.Wrap(errors.New(resp.Error), "failed to request")
		}
		channels = append(channels, resp.Channels...)
		if len(resp.Channels) == 0 || resp.ResponseMetadata.NextCursor == "" {
			break
		}
		cursor = resp.ResponseMetadata.NextCursor
	}
	return channels, nil
}

func fetchMessages(channelID string, since time.Time, token string) (messages []Message, err error) {
	messages = make([]Message, 0, 1000)
	u, err := url.Parse("https://slack.com/api/channels.history")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	v := u.Query()
	v.Set("token", token)
	v.Set("channel", channelID)
	v.Set("count", "1000")

	oldest := strconv.Itoa(int(since.Unix()))
	for {
		v.Set("oldest", oldest)

		u.RawQuery = v.Encode()

		resp := ChannelHistoryResponse{}
		if err = httpGet(u, &resp); err != nil {
			return nil, errors.Wrap(err, "failed to http get")
		}
		if !resp.Ok {
			return nil, errors.Wrap(errors.New(resp.Error), "failed to request")
		}
		messages = append(messages, resp.Messages...)
		log.Printf("messages count: %d", len(messages))
		if len(resp.Messages) == 0 || !resp.HasMore {
			break
		}
		oldest = resp.Messages[0].Ts
	}
	return messages, nil
}

func httpGet(u *url.URL, body interface{}) error {
	resp, err := http.Get(u.String())
	if err != nil {
		return errors.Wrap(err, "failed to http get")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		log.Printf("Retry-After: %s", retryAfter)
		a, err := strconv.Atoi(retryAfter)
		if err != nil {
			log.Printf("failed to parse Retry-After: %v", err)
			a = 1
		}
		time.Sleep(time.Duration(a) * time.Second)

		return httpGet(u, body)
	}
	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		return errors.Wrap(err, "failed to decode body")
	}
	return nil
}
