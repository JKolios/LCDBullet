package wsclient

import (
	"log"
	"fmt"
	"net/http"
	"html/template"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/gorilla/websocket")


type WebsocketConsumer struct {
	WSClientHost string
	WSClientEndpoint      string
	WSClientListenAddress string
}

var clientTemplate *template.Template
var host string
var upgrader = websocket.Upgrader{}
var wsContent = make(chan string)

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


func (consumer *WebsocketConsumer) Initialize(config conf.Configuration) {
	// Config Parsing
	consumer.WSClientHost = config.WSClientHost
	host = config.WSClientHost
	consumer.WSClientEndpoint = config.WSClientEndpoint
	consumer.WSClientListenAddress = config.WSClientListenAddress

	clientTemplate = template.Must(template.New("wsclient").Parse(clientTemplateStr))
	log.Println("Websocket Consumer: initialized, template loaded")
}

func (consumer *WebsocketConsumer) Register(EventInput chan events.Event) {

	// Input Monitor Goroutine Startup
	go monitorInput(EventInput)

	//Websocket Endpoint Startup
	http.HandleFunc("/dataSource", WSEndpointHandler)
	http.HandleFunc(consumer.WSClientEndpoint, ClientHandler)

	go http.ListenAndServe(consumer.WSClientListenAddress, nil)
	log.Println("Websocket Consumer: started, listening at " + consumer.WSClientHost + consumer.WSClientEndpoint)
}


func monitorInput(inputChan chan events.Event) {
	var incomingEvent events.Event

	for {
		incomingEvent = <- inputChan
		wsContent <- fmt.Sprintf("%s:%s\n", incomingEvent.Type, incomingEvent.Payload.(string))

	}

}

func WSEndpointHandler(w http.ResponseWriter, req *http.Request) {
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


func ClientHandler(w http.ResponseWriter, req *http.Request) {
	clientTemplate.Execute(w, "ws://" + host + "/dataSource")
}