function write(parent, message) {
    let item = document.createElement("div");
    for (let i = 0; i < message.body.length; i++) {
        let item = document.createElement("div");
        item.innerText = message.body[i];
        let doScroll = parent.scrollTop > parent.scrollHeight - parent.clientHeight - 1;
        parent.appendChild(item);
        if (doScroll) {
            parent.scrollTop = parent.scrollHeight - parent.clientHeight;
        }
    }
}

window.onload = function () {
    let conn;
    let log = document.getElementById("log");
    let cfgfile = document.getElementById("cfgfile");
    let vars = document.getElementById("env");
    let clear = document.getElementById("clear")

    // clear terminal window
    clear.onclick = function() { log.innerHTML = "" }

    if (window["WebSocket"]) {
        // creates a new WebSocket connection
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onopen = function(event) {
            write(log, { type:2, body: ["Connection established."]});
        }
        conn.onclose = function (evt) {
            write(log, { type:2, body: ["Connection closed."]});
        };
        conn.onmessage = function (evt) {
            // parse the message string into a JSON object
            let message = JSON.parse(evt.data);

            // the message is an event
            if (message.type == 0) {
                write(log, message)
            }
            // the message contain one or more configuration files
            if (message.type == 1) {
                write(cfgfile, message)
            }
            // the message has environment variables
            if (message.type == 2) {
                write(vars, message)
            }
        };
        conn.onerror = function (evt) {
            write(log, { type:2, body: ["Connection error."]});
        };
    } else {
        write(log, { type:2, body: ["Your browser does not support WebSockets."]});
    }
};