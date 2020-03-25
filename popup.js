$(document).ready(function () {

    InitializeEventHandlers();
});

function InitializeEventHandlers() {

    $('#sergingSwitch').on("click", function () {
        StartSerging();
    });

    $("#connectedChannels tbody").on("click", 'input', '#connectedChannels', function (event) {
        let channame = event.target.textContent;
        if (event.target.checked)
        {
            SendChannelConnection(channame);
        }
        else
        {
            SendChannelDisconnection(channame);
        }
    });
}


function StartSerging() {
    chrome.runtime.sendMessage({ message: "connect_to_websocket" }, function (response) {
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
    console.log(tabs.length)
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
                    checkCell.appendChild(chkElement);
                }
            });
        }
    }
});