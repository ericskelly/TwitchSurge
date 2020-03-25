package server

import (
	"fmt"
	"strings"
	"sync"
)

var globalSubscribedChannels map[string]*ChannelConnection = make(map[string]*ChannelConnection)
var globalSubscribedLock sync.RWMutex

var baseFFZUrl string = "https://api.frankerfacez.com/v1/room/"
var baseBTTVUrl string = "https://api.betterttv.net/2/channels/"

//SubscribeToChannel - subcribe a user to a channel
func SubscribeToChannel(channel string, client *ClientInfo) {
	currentChannel := globalSubscribedChannels[channel]
	if currentChannel == nil {
		globalSubscribedLock.Lock()
		currentChannel = &ChannelConnection{}
		currentChannel.members = []*ClientInfo{client}
		currentChannel.channelName = channel
		globalSubscribedChannels[channel] = currentChannel
		getChannelFFZEmotes(currentChannel)
		getBTTVEmotes(currentChannel)
		IrcConnection(channel, currentChannel)
		globalSubscribedLock.Unlock()
	} else {
		currentChannel.Lock()
		addToSliceClientInfo(&currentChannel.members, client)
		currentChannel.Unlock()
	}
}

//UnsubscribeFromChannel - unsubscribe a user from a channel
func UnsubscribeFromChannel(channel string, client *ClientInfo) {

}

func getChannelFFZEmotes(channelConnection *ChannelConnection) {
	channelURL := baseFFZUrl + strings.ToLower(channelConnection.channelName)
	response := FFZRoomResponse{}
	err := GetJSON(channelURL, &response)
	if err != nil {
		return
	}

	var ffzChannelEmotes map[string]string = make(map[string]string)
	for key := range response.Sets {
		for _, emote := range response.Sets[key].Emoticons {
			ffzChannelEmotes[emote.Name] = emote.Name
		}
	}

	channelConnection.currentFFZEmoteNames = ffzChannelEmotes
}

func getBTTVEmotes(ChannelConnection *ChannelConnection) {
	channelURL := baseBTTVUrl + strings.ToLower(ChannelConnection.channelName)
	response := BTTVResponse{}
	err := GetJSON(channelURL, &response)
	if err != nil {
		return
	}

	var bttvChannelEmotes map[string]string = make(map[string]string)
	for _, emote := range response.Emotes {
		bttvChannelEmotes[emote.Code] = emote.Code
	}

	ChannelConnection.currentBTTVEmoteNames = bttvChannelEmotes
	fmt.Println(ChannelConnection.currentBTTVEmoteNames)
}

func addToSliceClientInfo(slice *[]*ClientInfo, clientinfo *ClientInfo) bool {
	newslice := *slice
	for _, v := range newslice {
		if v == clientinfo {
			return false
		}
	}

	newslice = append(newslice, clientinfo)
	*slice = newslice
	return true
}
