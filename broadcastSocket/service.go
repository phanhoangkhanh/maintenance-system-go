package broadcastsocket

import (
	"log"
	"maintenance-system-go/config"
	"maintenance-system-go/tradingSocket"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//PROCESS CLOSE CONNECTION WITH CLIENT
// 1- FULL MSG -> SEND ALARM TO CHANNEL signalDone
// 				-> ALL GOROUTINES RELATED TO CLIENT EXIT 'FOR ROUTINE' AND CALL STOP CLIENT + REMOVE LIST CLIENTS
// 2- CLIENT DISCONNECT -> READ OR WRITE FAIL
// 				-> RETURN ROUTINE AND CALL TO CLOSE CLIENT CONNECTION

type Service struct {
	Config        *config.Config
	Conn          *websocket.Conn
	lock          sync.Mutex
	BridgeChannel chan tradingSocket.MarketPriceInfo
}

//Help socket upgrade with allowed origins in the CORS configuration
func newWebSocketUpgrader(allowedOrigins []string) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "" || slices.Contains(allowedOrigins, origin)
		},
	}
}

type Client struct {
	conn         *websocket.Conn
	msgSend      chan tradingSocket.MarketPriceInfo //private channel to store and forward market price info to this client
	lock         sync.Mutex
	isReady      bool                 // flag to indicate if the client is ready
	signalDone   chan struct{}        // channel to signal when the client is done-> close everything related
	callCloseNow func(client *Client) // callback to close the client connection
}

type ClientList struct {
	clients []*Client
	lock    sync.Mutex
}

// use Struct with lock to make sure no data-race occurs when accessing the clients list from many goroutines
var clientList = &ClientList{
	clients: []*Client{},
	lock:    sync.Mutex{},
}

func (s *Service) CreatAndBroadcastSocket() {

	//Let some URI request origins be allowed based on the CORS configuration
	var allowedOrigins []string
	if s.Config != nil {
		allowedOrigins = s.Config.CORS.AllowedOrigins
	}
	upgrader := newWebSocketUpgrader(allowedOrigins)

	//1 - Multiply msg from origin Bridge to all connected clients
	go multipleMsgToClients(s.BridgeChannel)

	//2 - This HandleFunc will be called each time has one request call to '/ws' socket
	//each time user connect. ->call 1 time and create new 1 client
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//Upgrate to socket -> handshake with the client
		webConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Fail to upgrate http request to socket, %s\n", err)
			return
		}

		//0 - init the only one callback func to close each client
		var onceClose sync.Once
		var callCloseNow = func(client *Client) {
			onceClose.Do(func() {
				// no need clock client -> close the connection directly
				client.conn.Close()
				close(client.signalDone)

				removeClientFromList(client)
				log.Println("remove client done....")
			})
		}

		//create new Client struct for each request connection
		newClient := &Client{
			conn:         webConn,
			msgSend:      make(chan tradingSocket.MarketPriceInfo, 1000),
			lock:         sync.Mutex{},
			isReady:      false,
			signalDone:   make(chan struct{}), //the signal channel to indicate done, close routiune read and write
			callCloseNow: callCloseNow,
		}
		// add new client to the list of clients
		clientList.lock.Lock()
		clientList.clients = append(clientList.clients, newClient)
		clientList.lock.Unlock()

		//Keep reading from and writing to the client in separate goroutines
		go ReadFromClient(newClient)
		go WriteToClient(newClient)
	})
	//This is infinite loop to keep HandleFunc live and wait
	http.ListenAndServe(":8000", nil)
}

func multipleMsgToClients(Broadcast chan tradingSocket.MarketPriceInfo) {
	//get per message from channel
	// gorountine block -> keep waiting for msg -> will end when channel is closed and all messages are processed
	// Channel empty but not closed, will keep waiting for new messages
	for msg := range Broadcast {
		// 1 msg from the broadcast channel, multiple clients receive it
		clientList.lock.Lock()
		snapshot := append([]*Client(nil), clientList.clients...) // use a snapshot to avoid holding the lock while sending messages to clients
		clientList.lock.Unlock()
		for _, client := range snapshot {
			//This simple way will let client msgSend be fulled even 1000 buffer and block everything
			//client.msgSend <- msg
			select {
			// routine block -> wait client.signalDone or client.msgSend
			// 		-> cannot do that -> default -> close client
			case client.msgSend <- msg: //if can send more msg to client , just do it
			case <-client.signalDone: // if client is closed, skip sending
				continue
			default:
				client.callCloseNow(client) // Full buffer, force close the client connection
				continue
			}
		}
	}
}

// When client disconnect -> ReadFromClient recognizes it and removes the client from the list
func ReadFromClient(client *Client) {
	defer client.callCloseNow(client) // close connect of this client
	for {
		//check if the client is done
		select {
		case <-client.signalDone:
			return //close the goroutine ReadFromClient when client is done
		default:
			var msg map[string]string
			//also goroutine block -> will wait for client to send a message
			err := client.conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("Error reading from client: %s\n", err)
				return
			}
			// Forward the message to the client's private channel
			if msg["login"] == "KHANH SECRET" {
				messageSuccess := "Login successfully"
				client.lock.Lock()

				if err := client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
					log.Printf("Error setting client write deadline: %s\n", err)
					client.lock.Unlock()
					return
				}
				err := client.conn.WriteJSON(map[string]string{"message": messageSuccess})
				if err != nil {
					log.Printf("Error writing to client: %s\n", err)
					client.lock.Unlock()
					return
				}
				client.isReady = true
				client.lock.Unlock()
			}
		}
	}
}

func WriteToClient(client *Client) {
	defer client.callCloseNow(client) // close connect of this client
	for {
		// routine block -> wait client.signalDone or client.msgSend (no default so keep waiting 2 channels)
		select {
		case <-client.signalDone:
			return //close the goroutine WriteToClient if client is done
		case msg := <-client.msgSend:
			client.lock.Lock()
			if !client.isReady {
				client.lock.Unlock()
				continue
			}
			//Set deadline for writing to the client to avoid blocking indefinitely
			if err := client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Error setting client write deadline: %s\n", err)
				client.lock.Unlock()
				return
			}
			//also goroutine block
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing to client: %s\n", err)
				client.lock.Unlock()
				return
			}
			client.lock.Unlock()

		}

	}

}

func removeClientFromList(client *Client) {
	clientList.lock.Lock()
	for i, c := range clientList.clients {
		if c == client {
			clientList.clients = append(clientList.clients[:i], clientList.clients[i+1:]...)
			break
		}
	}
	clientList.lock.Unlock()
}
