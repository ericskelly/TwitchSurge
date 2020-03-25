package server

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"strings"
	"time"
)

// IrcConnection irc connection client
func IrcConnection(channel string, channelConnection *ChannelConnection) {
	runIrcConnection(channel, channelConnection)
}

func runIrcConnection(channel string, channelConnection *ChannelConnection) {
	channelName := "#" + strings.ToLower(channel)
	c := irc.SimpleClient("ttvserge")
	c.Config().Server = "irc.chat.twitch.tv:6667"
	c.Config().Pass = Configuration.TwitchAuth

	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		fmt.Println("connected")
		conn.Join(channelName)
	})

	c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		msg := line.Args[1]
		channelConnection.weightedMessagesPerInterval += determineWeightedChatMessage(line, channelConnection)
		fmt.Println(channelName, line.Nick, msg)
	})

	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		fmt.Println("disconnected")
	})

	if err := c.Connect(); err != nil {
		fmt.Printf("Irc Connection error: %s\n", err.Error())
	}

	go runTicker(channelConnection)

}

func determineWeightedChatMessage(line *irc.Line, channelConnection *ChannelConnection) float64 {
	weightedMessage := 1.0
	msg := line.Args[1]
	var words []string = strings.Split(msg, " ")
	fmt.Println(words)
	for _, val := range words {
		if _, ok := channelConnection.currentFFZEmoteNames[val]; ok {
			weightedMessage = weightedMessage * 2
			break
		} else if _, ok := channelConnection.currentBTTVEmoteNames[val]; ok {
			weightedMessage = weightedMessage * 2
			break
		}

	}

	return weightedMessage
}

func runTicker(channelConnection *ChannelConnection) {
	channelConnection.channelTicker = time.NewTicker(10 * time.Second)
	averageMessageCounter := 0
	messagesPerMinute := 0.0
	for {
		select {
		case <-channelConnection.channelTicker.C:
			messagesPerMinute += channelConnection.weightedMessagesPerInterval
			fmt.Println(messagesPerMinute)
			averageMessageCounter++
			if averageMessageCounter == 6 {
				channelConnection.averageMessagesPerInterval = messagesPerMinute / 6.0
				messagesPerMinute = 0.0
				averageMessageCounter = 0
			}
			if channelConnection.averageMessagesPerInterval > 0 {
				checkForSurge(channelConnection)
			}
			channelConnection.weightedMessagesPerInterval = 0
		}
	}
}

func checkForSurge(channelConnection *ChannelConnection) {
	fmt.Println(channelConnection.weightedMessagesPerInterval, channelConnection.averageMessagesPerInterval, channelConnection.averageMessagesPerInterval*1.5)
	if channelConnection.weightedMessagesPerInterval > (channelConnection.averageMessagesPerInterval * 1.5) {
		sendSurgeMessages(channelConnection)

	}
}

func sendSurgeMessages(channelConnection *ChannelConnection) {
	surgeMessage := SurgeMessage{ChannelName: channelConnection.channelName}
	for _, client := range channelConnection.members {
		client.surgeChannel <- surgeMessage
	}
}
