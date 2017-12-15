package main

type ChannelListResponse struct {
	Channels         []Channel        `json:"channels"`
	Ok               bool             `json:"ok"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
	Error            string           `json:"error"`
}

type Channel struct {
	Created        int64         `json:"created"`
	Creator        string        `json:"creator"`
	ID             string        `json:"id"`
	IsArchived     bool          `json:"is_archived"`
	IsChannel      bool          `json:"is_channel"`
	IsGeneral      bool          `json:"is_general"`
	IsMember       bool          `json:"is_member"`
	IsMpim         bool          `json:"is_mpim"`
	IsOrgShared    bool          `json:"is_org_shared"`
	IsPrivate      bool          `json:"is_private"`
	IsShared       bool          `json:"is_shared"`
	Members        []string      `json:"members"`
	Name           string        `json:"name"`
	NameNormalized string        `json:"name_normalized"`
	NumMembers     int64         `json:"num_members"`
	PreviousNames  []interface{} `json:"previous_names"`
	Purpose        Info          `json:"purpose"`
	Topic          Info          `json:"topic"`
}

type Info struct {
	Creator string `json:"creator"`
	LastSet int64  `json:"last_set"`
	Value   string `json:"value"`
}

type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}
