package wsclient

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/JKolios/EventsToGo/consumers"
	"github.com/JKolios/EventsToGo/events"
	"github.com/gorilla/websocket"
)

const clientTemplateStr string = `
<!DOCTYPE html>
<head>
    <meta charset="utf-8">
    <script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("Connection Opened");
        }
        ws.onclose = function(evt) {
            print("Connection Closed");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print(evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
	<table>
		<tr>
			<td valign="top" width="50%">
				<p>Click "Open" to create a connection to the server and "Close" to close the connection.
					<form>
						<button id="open">Open</button>
						<button id="close">Close</button>
					</form>
	</table>
	<div id ="output"></div>
</body>
</html>`

func WSEndpointClosure(wsContent chan string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		upgrader := &websocket.Upgrader{}
		wsConn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Websocket Consumer: upgrade error:", err)
			return
		}
		defer wsConn.Close()

		for {
			err = wsConn.WriteMessage(websocket.TextMessage, []byte(<-wsContent))
			if err != nil {
				log.Println("Websocket Consumer: write error:", err)
				break
			}
		}
	}
}

func ClientHandlerClosure(template *template.Template, host string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		template.Execute(w, "ws://"+host+"/dataSource")
	}
}

func RunFunction(consumer *consumers.GenericConsumer, incomingEvent events.Event) {
	consumer.RuntimeObjects["wsContentChan"].(chan string) <- fmt.Sprintf("%s:%s\n", incomingEvent.Type, incomingEvent.Payload.(string))
}

func SetupFunction(consumer *consumers.GenericConsumer, config map[string]interface{}) {

	print(config)
	consumer.RuntimeObjects["template"] = template.Must(template.New("wsclient").Parse(clientTemplateStr))
	consumer.RuntimeObjects["wsContentChan"] = make(chan string)

	//Websocket Endpoint Startup
	http.HandleFunc("/dataSource", WSEndpointClosure(consumer.RuntimeObjects["wsContentChan"].(chan string)))
	http.HandleFunc(config["WSClientEndpoint"].(string), ClientHandlerClosure(consumer.RuntimeObjects["template"].(*template.Template), config["WSClientHost"].(string)))

	go http.ListenAndServe(config["WSClientListenAddress"].(string), nil)
	log.Println("Websocket Endpoint: started, listening at " + config["WSClientHost"].(string) + config["WSClientEndpoint"].(string))

}
