var websocket;
chrome.runtime.onInstalled.addListener(function () {

  chrome.declarativeContent.onPageChanged.removeRules(undefined, function () {
    chrome.declarativeContent.onPageChanged.addRules([{
      conditions: [new chrome.declarativeContent.PageStateMatcher({
        pageUrl: { hostEquals: 'www.twitch.tv', schemes: ['https'] },
      })],
      actions: [new chrome.declarativeContent.ShowPageAction()]
    }]);
  });
});

chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
  if (request.message === "connect_to_websocket")
  {
    this.createWebSocketConnection();
    sendResponse({ message: "received" });
  }
  else if (request.message === "disconnect_from_websocket")
  {
    this.closeWebSocketConnection()
    sendResponse({ message: "received" });
  }
  else if (request.message === "send_channel_sub")
  {
    this.sendChannelSubOrUnsub(request.channelname, request.type);
    sendResponse({ message: "received" });
  }
  else if (request.message === "send_channel_unsub")
  {
    this.sendChannelSubOrUnsub(request.channelname, request.type);
    sendResponse({ message: "received" });
  }
});

function createWebSocketConnection() {
  if ('WebSocket' in window)
  {
    connect('ws://localhost:5000/ws');
  }
}

//Make a websocket connection with the server.
function connect(host) {
  if (!websocket)
  {
    websocket = new WebSocket(host);
  }

  websocket.onopen = function () {

    websocket.send(JSON.stringify({ channelname: 'open' }));
    chrome.storage.local.set({ 'websocketConn': true }, function () {
      console.log("websocket connection saved")
    });
  };

  websocket.onmessage = function (event) {
    var received_msg = JSON.parse(event.data);
    console.log(received_msg)
    var channelname = received_msg.ChannelName
    var demoNotificationOptions = {
      type: "basic",
      message: channelname + " Is Surging!",
      title: "Twitch Surge",
      iconUrl: "images/twitch.png"
    }
    chrome.notifications.create("", demoNotificationOptions);

  };

  websocket.onclose = function () {
    closeWebSocketConnection()
  };
}

function sendChannelSubOrUnsub(channelname, type) {
  if (websocket)
  {
    websocket.send(JSON.stringify({ channelName: channelname, type: type }));
  }
}

function closeWebSocketConnection() {
  console.log(websocket);
  if (websocket != null || websocket != undefined)
  {
    websocket.close();
    websocket = undefined;
    chrome.storage.local.set({ 'websocketConn': false }, function () {
      console.log("websocket connection closed")
    });

    chrome.storage.local.remove('checkedSettings');
  }
}
