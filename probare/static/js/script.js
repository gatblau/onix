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

function getBannerClass(type) {
    switch (type) {
        case "":
            return "pf-c-banner banner-font-size";
        case "info":
            return "pf-c-banner pf-m-info banner-font-size";
        case "danger":
            return "pf-c-banner pf-m-danger banner-font-size";
        case "success":
            return "pf-c-banner pf-m-success banner-font-size";
        case "warning":
            return "pf-c-banner pf-m-warning banner-font-size";
    }
}

window.onload = function () {
    let conn;
    let log = document.getElementById("log");
    let cfgfile = document.getElementById("cfgfile");
    let vars = document.getElementById("env");
    let clear = document.getElementById("clear")
    let banner = document.getElementById("banner")

    // clear terminal window
    clear.onclick = function() { log.innerHTML = "" }

    if (window["WebSocket"]) {
        // creates a new WebSocket connection
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onopen = function(event) {
            write(log, { type:2, body: ["connected to the server"]});
        }
        conn.onclose = function (evt) {
            write(log, { type:2, body: ["lost connection to the server"]});
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
                cfgfile.innerHTML = ""
                write(cfgfile, message)
            }
            // the message has environment variables
            if (message.type == 2) {
                vars.innerHTML = ""
                write(vars, message)
            }
            // control message
            if (message.type == 3) {
                // set the banner class
                banner.className = getBannerClass(message.body[0]);
                // set the banner message
                banner.innerHTML = message.body[1];
            }
        };
        conn.onerror = function (evt) {
            write(log, { type:2, body: ["Connection error."]});
        };
    } else {
        write(log, { type:2, body: ["Your browser does not support WebSockets."]});
    }
};