package live

type LiveMediaInfo struct {
	HasAudio     bool   `json:"hasAudio"`
	OnlineStatus bool   `json:"onlineStatus"`
	AnchorName   string `json:"anchorName"`
}
