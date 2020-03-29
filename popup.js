$(document).ready(function () {

    CheckWebsocketConn();
    InitializeEventHandlers();
});

function CheckWebsocketConn() {
    chrome.storage.local.get(['websocketConn'], function (result) {
        console.log(result.websocketConn);
        if (result.websocketConn && result.websocketConn == true)
        {
            var surgingSwitch = document.getElementById("surgingSwitch");
            surgingSwitch.checked = true;
        }
    });
}

function InitializeEventHandlers() {

    $('#surgingSwitch').on("click", function () {
        console.log($(this).prop('checked'));
        if ($(this).prop('checked'))
        {
            StartSurging()
        }
        else
        {
            StopSurging();
        }
    });

    $("#connectedChannels tbody").on("click", 'input', '#connectedChannels', function (event) {
        let channame = event.target.textContent;
        let chanNameCheckedList = [];
        chrome.storage.local.get(['checkedSettings'], function (result) {
            chanNameCheckedList = result.checkedSettings ? result.checkedSettings : [];
            if (event.target.checked)
            {
                chanNameCheckedList.push(channame);
                SendChannelConnection(channame);
            }
            else
            {
                for (let i = 0; i < chanNameCheckedList.length; ++i)
                {
                    if (chanNameCheckedList[i] == channame)
                    {
                        chanNameCheckedList.splice(i, 1);
                    }
                }
                SendChannelDisconnection(channame);
            }
            chrome.storage.local.set({ 'checkedSettings': chanNameCheckedList }, function () {
                console.log("check saved");
            });
        });
    });
}

function StartSurging() {
    chrome.runtime.sendMessage({ message: "connect_to_websocket" }, function (response) {
        console.log(response.message);
    });
}

function StopSurging() {
    chrome.runtime.sendMessage({ message: "disconnect_from_websocket" }, function (response) {
        console.log(response.message);

    });
}

function SendChannelConnection(channelname) {
    chrome.runtime.sendMessage({ message: "send_channel_sub", channelname: channelname, type: "sub" }, function (response) {
        console.log(response.message);
    });
}

function SendChannelDisconnection(channelname) {
    chrome.runtime.sendMessage({ message: "send_channel_unsub", channelname: channelname, type: "unsub" }, function (response) {
        console.log(response.message);
    });
}

chrome.tabs.query({}, function (tabs) {
    chrome.storage.local.get(['checkedSettings'], function (result) {
        var resultsArray = result.checkedSettings;
        for (i = 0; i < tabs.length; ++i)
        {
            if (tabs[i].url.includes("twitch.tv"))
            {
                chrome.tabs.sendMessage(tabs[i].id, { message: "get_channel_name" }, function (response) {
                    if (response)
                    {
                        console.log(response);
                        var table = document.getElementById("connectedChannels");
                        var row = table.insertRow();
                        var channelCell = row.insertCell(0);
                        var checkCell = row.insertCell(1);
                        channelCell.innerHTML = response.channelname;
                        var chkElement = document.createElement("input");
                        chkElement.type = "checkbox";
                        chkElement.innerText = response.channelname;
                        if (resultsArray && resultsArray.includes(response.channelname))
                        {
                            chkElement.checked = true;
                        }
                        checkCell.appendChild(chkElement);
                    }
                });
            }
        }
    });
});