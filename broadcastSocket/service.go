package broadcastsocket

import (
	"log"
	"maintenance-system-go/config"
	"maintenance-system-go/tradingSocket"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//PROCESS CLOSE CONNECTION WITH CLIENT
// 1- FULL MSG -> SEND ALARM TO CHANNEL signalDone
// 				-> ALL GOROUTINES RELATED TO CLIENT EXIT 'FOR ROUTINE' AND CALL STOP CLIENT + REMOVE LIST CLIENTS
// 2- CLIENT DISCONNECT -> READ OR WRITE FAIL
// 				-> RETURN ROUTINE AND CALL TO CLOSE CLIENT CONNECTION

type Service struct {
	Config *config.Config
	Conn           *websocket.Conn
	lock           sync.Mutex
	BridgeChannel  chan tradingSocket.MarketPriceInfo
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Client struct {
	conn       *websocket.Conn
	msgSend      chan tradingSocket.MarketPriceInfo //private channel to store and forward market price info to this client
	lock         sync.Mutex
	isReady      bool // flag to indicate if the client is ready
	signalDone   chan struct{} // channel to signal when the client is done-> close everything related
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

func (s *Service) CreatAndBroadcastSocket()  {


	//1 - Multiply msg from origin Bridge to all connected clients 
	go multipleMsgToClients(s.BridgeChannel)


	//2 - This HandleFunc will be called each time has one request call to '/ws' socket
	//each time user connect. ->call 1 time and create new 1 client
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//Upgrate to socket
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
				log.Println("remove client now....")
			})
		}

		//create new Client struct for each request connection
		newClient := &Client{
			conn:       webConn,
			msgSend:       make(chan tradingSocket.MarketPriceInfo, 1000),
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
		go ReadFromClient(newClient )
		go WriteToClient(newClient )
	})
	//This is infinite loop to keep HandleFunc live and wait
	http.ListenAndServe(":8000", nil)
}

func multipleMsgToClients(Broadcast chan tradingSocket.MarketPriceInfo) {
	//get per message from channel
	for msg := range Broadcast {
		// 1 msg from the broadcast channel, multiple clients receive it
		clientList.lock.Lock()
		for _, client := range clientList.clients {
			//This simple way will let client msgSend be fulled even 1000 buffer and block everything
			//client.msgSend <- msg
			select {
				case client.msgSend <- msg: //if can send more msg to client , just do it
				case <-client.signalDone: // if client is closed, skip sending
					continue
				default:
					client.callCloseNow(client) // Full buffer, force close the client connection
					continue
			}
		}
		clientList.lock.Unlock()
	}
}

//When client disconnect -> ReadFromClient recognizes it and removes the client from the list
func ReadFromClient(client *Client) {
	defer client.callCloseNow(client) // close connect of this client
	for {
		//check if the client is done
		select {
			case <-client.signalDone:
				return //close the goroutine ReadFromClient when client is done
			default:
				var msg map[string]string
				err := client.conn.ReadJSON(&msg)
				if err != nil {
					log.Printf("Error reading from client: %s\n", err)
					return
				}
				// Forward the message to the client's private channel
				if msg["login"] == "KHANH SECRET" {
					messageSuccess := "Login successfully"
					client.lock.Lock()
					client.conn.WriteJSON(map[string]string{"message": messageSuccess})
					client.isReady = true
					client.lock.Unlock()
				}
		}
	}
}

func WriteToClient(client *Client) {
	defer client.callCloseNow(client) // close connect of this client
	for{
		select{
			case <-client.signalDone:
				return	//close the goroutine WriteToClient if client is done
			case msg := <-client.msgSend: 
				client.lock.Lock()
				if !client.isReady {
					client.lock.Unlock()
					continue
				}
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


func removeClientFromList( client *Client) {
	clientList.lock.Lock()
	for i, c := range clientList.clients {
		if c == client {
			clientList.clients = append(clientList.clients[:i], clientList.clients[i+1:]...)
			break
		}
	}
	clientList.lock.Unlock()
}
