package server

import (
	//"fmt"
	"strings"
	"sync"
)

var globalSubscribedChannels map[string]*ChannelConnection = make(map[string]*ChannelConnection)
var globalSubscribedLock sync.RWMutex

var baseFFZUrl string = "https://api.frankerfacez.com/v1/room/"
var baseBTTVUrl string = "https://api.betterttv.net/2/channels/"

//SubscribeToChannel - subcribe a user to a channel
func SubscribeToChannel(channel string, client *ClientInfo) {
	globalSubscribedLock.RLock()
	currentChannel := globalSubscribedChannels[channel]
	if currentChannel == nil {
		globalSubscribedLock.Lock()
		currentChannel = &ChannelConnection{}
		currentChannel.members = []*ClientInfo{client}
		currentChannel.channelName = channel
		globalSubscribedChannels[channel] = currentChannel
		getChannelEmotesConcurrentWait(currentChannel)
		IrcConnection(channel, currentChannel)
		globalSubscribedLock.Unlock()
	} else {
		currentChannel.Lock()
		addToSliceClientInfo(&currentChannel.members, client)
		currentChannel.Unlock()
	}
	globalSubscribedLock.RUnlock()
}

//UnsubscribeFromChannel - unsubscribe a user from a channel
func UnsubscribeFromChannel(channel string, client *ClientInfo) {
	globalSubscribedLock.RLock()
	channelConn := globalSubscribedChannels[channel]
	if channelConn != nil {
		globalSubscribedLock.Lock()
		removeFromSliceClientInfo(&channelConn.members, client)
		if len(channelConn.members) == 0 {
			IrcDisconnection(channelConn)
			delete(globalSubscribedChannels, channel)
		}
		globalSubscribedLock.Unlock()
	}
	globalSubscribedLock.RUnlock()
}

func getChannelEmotesConcurrentWait(channelConnection *ChannelConnection) {
	var wg sync.WaitGroup
	wg.Add(2)
	go getChannelFFZEmotes(channelConnection, &wg)
	go getBTTVEmotes(channelConnection, &wg)
	wg.Wait()
}

func getChannelFFZEmotes(channelConnection *ChannelConnection, wg *sync.WaitGroup) {
	defer wg.Done()
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

func getBTTVEmotes(channelConnection *ChannelConnection, wg *sync.WaitGroup) {
	defer wg.Done()
	channelURL := baseBTTVUrl + strings.ToLower(channelConnection.channelName)
	response := BTTVResponse{}
	err := GetJSON(channelURL, &response)
	if err != nil {
		return
	}

	var bttvChannelEmotes map[string]string = make(map[string]string)
	for _, emote := range response.Emotes {
		bttvChannelEmotes[emote.Code] = emote.Code
	}

	channelConnection.currentBTTVEmoteNames = bttvChannelEmotes
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

func removeFromSliceClientInfo(slice *[]*ClientInfo, clientinfo *ClientInfo) bool {
	newslice := *slice
	removed := false
	for i, v := range newslice {
		if v == clientinfo {
			newslice = append(newslice[:i], newslice[i+1:]...)
			removed = true
		}
	}
	*slice = newslice
	return removed
}
