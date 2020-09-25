window.onload = function () {
    let conn;
    let log = document.getElementById("log");
    let clear = document.getElementById("clear")

    clear.onclick = function() {
        log.innerHTML = ""
    }

    function appendLog(item) {
        let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    if (window["WebSocket"]) {
        // creates a new WebSocket connection
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onopen = function(event) {
            let item = document.createElement("div");
            item.innerHTML = "<b>Connection established.</b>";
            appendLog(item);
        }
        conn.onclose = function (evt) {
            let item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            let messages = evt.data.split('\n');
            for (let i = 0; i < messages.length; i++) {
                let item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
        conn.onerror = function (evt) {
            let item = document.createElement("div");
            item.innerHTML = "<b>Connection error.</b>";
            appendLog(item);
        };
    } else {
        let item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};