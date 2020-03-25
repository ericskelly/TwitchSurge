package server

import (
	"sync"
	"time"
)

//SurgeMessage - Contains information for the surge message
type SurgeMessage struct {
	ChannelName string
}

//SurgeConnect - Contains information to subscribe or unsubscribe from a channel
type SurgeConnect struct {
	ChannelName string
	Type        string
}

//ClientInfo - info for the connected client
type ClientInfo struct {
	twitchusername    string
	connectedChannels []string
	surgeChannel      chan<- SurgeMessage
}

//Config - configuration
type Config struct {
	OriginsAllowed []string
	TwitchAuth     string
}

//ChannelConnection - contains information for a channel connection
type ChannelConnection struct {
	channelName                 string
	averageMessagesPerInterval  float64
	weightedMessagesPerInterval float64
	isSurging                   bool
	pastMessagesBuffer          []string
	channelTicker               *time.Ticker
	members                     []*ClientInfo
	currentFFZEmoteNames        map[string]string
	currentBTTVEmoteNames       map[string]string
	sync.RWMutex
}

//FFZRoomResponse - response type for frankerfacez room information
type FFZRoomResponse struct {
	Room FFZRoom
	Sets map[string]FFZEmoteSet
}

//FFZRoom - FrankerFaceZ room data
type FFZRoom struct {
	RoomName string `json:"display_name"`
}

//FFZEmoteSet - Emote set data
type FFZEmoteSet struct {
	Emoticons []FFZEmote
}

//FFZEmote - Data for single frankerfacez emote
type FFZEmote struct {
	Name string
}

//BTTVResponse - response type for betterttv channel information
type BTTVResponse struct {
	Emotes []BTTVEmote
}

//BTTVEmote - Data for single betterttv emote
type BTTVEmote struct {
	Code string
}
