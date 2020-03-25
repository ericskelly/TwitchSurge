window.onload = function () {
    this.GetChannelName()
}

function GetChannelName() {
    console.log("-------------------TwitchSurge-----------------------")
    var channelElement = document.getElementsByClassName("tw-c-text-inherit tw-font-size-5 tw-white-space-nowrap")[0];
    if (channelElement)
    {
        return channelElement.innerHTML;
    }

}

chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    var channelname;
    if (request.message === "get_channel_name")
    {
        channelname = this.GetChannelName();
        console.log(channelname)
        if (channelname)
        {
            sendResponse({ channelname: channelname })
        }
    }

});

