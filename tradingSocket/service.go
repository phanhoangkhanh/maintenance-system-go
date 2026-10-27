package tradingSocket

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"maintenance-system-go/config"
	"maintenance-system-go/database/redis"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Service struct {
	Redis          *redis.RedisClient
	Config         *config.Config
	Conn           *websocket.Conn
	lock           sync.Mutex
	IsReconnecting bool
	BridgeChannel chan MarketPriceInfo
}

func (ts *Service) ConnectWebsocket() {
	closeAllRoutine := make(chan struct{})
	keepGotMsg := make(chan struct{}, 1)
	var conn *websocket.Conn
	var reconnectOnce sync.Once
	reconnect := func() {
		// make sure only one reconnect is triggered at a time ConnectWebsocket called
		reconnectOnce.Do(func() {
			//close will init msg to all routine places: closeAllRoutine <-
			close(closeAllRoutine)
			ts.reconnectSocket()
		})
	}
	//STEP 1: connect URI socket
	wss_url := ts.Config.TradingWS.URI
	for {
		if newConn, err := connect(wss_url); err == nil {
			conn = newConn
			ts.lock.Lock()
			ts.Conn = conn
			ts.lock.Unlock()
			break
		} else {
			fmt.Println("Connection Trading - Wss failed, Retrying in 5 seconds...", err)
		}

		time.Sleep(5 * time.Second)
	}

	//2- SEND MSG KEY TO REGISTER AND GET INFO
	message := TradingKeyMessage{
		Key: TradingKeyInsideKeyMessage{
			DomainName:  "",
			Prefix:      "",
			MessageName: "subscribe",
			Suffix:      "v1",
			MessageType: "tick_price",
		},
		Payload: "miniTicker",
	}
	messageBytes, _ := json.Marshal(message)
	err := ts.SendMessage(messageBytes)
	if err != nil {
		//Re send after 5 seconds
		time.Sleep(5 * time.Second)
		ts.SendMessage(messageBytes)
	}

	//3- ROUTINE READ MSG
	go func() {
		for {
			//block and wait untill new message arrives
			_, message, err := conn.ReadMessage()
			if err != nil {
				//log.Printf("Connection lost: %v", err)
				log.Printf("Connection with Trading server lost: %v", err)
				reconnect()
				return
			}
			select {
			case <-closeAllRoutine:
				return
			case keepGotMsg <- struct{}{}:
			default:
			}

			var msg MessageTrading
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal trading message: %v", err)
				continue
			}

			// Decode base64 payload
			for i, s := range msg.Payload {
				//fmt.Println("ORIGIN....", string(s))
				decoded, err := base64.StdEncoding.DecodeString(s)
				if err != nil {
					log.Printf("payload[%d] not as base64 type: %v", i, err)
					continue
				}

				var price MarketPriceInfo
				if err := json.Unmarshal(decoded, &price); err != nil {
					log.Printf("Price Info Unmarshal error: %v", err)
					continue
				}

				//key := "S-SYSTEM-SYMBOL-" + price.Symbol
				//fmt.Printf("GOT: %s  %#v\n", key, price)
				// Send the price info to the bridge channel
				ts.BridgeChannel <- price
			}

			// Close for together with monitoring routine if waiting more than 70s
			select {
			case <-closeAllRoutine:
				return
			default:
			}
		}
	}()

	//4- Monitor keepGotMsg channel to detect if connection is alive
	go func() {
		timer := time.NewTimer(70 * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-keepGotMsg:
				//Reset timer every time a new message is received
				timer.Reset(70 * time.Second)

			case <-timer.C:
				log.Println("No message received for 70 seconds, reconnecting...")
				reconnect()
				return
			case <-closeAllRoutine: //got msg from listen routine
				log.Println("Connection lost, stopping monitoring routine...")
				return
			}

		}

	}()

}

func connect(url string) (*websocket.Conn, error) {

	dialer := websocket.Dialer{
		NetDial:          nil,
		HandshakeTimeout: 0,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification
		},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to Trading-WebSocket server")
	return conn, nil
}

func (ts *Service) SendMessage(message []byte) error {
	//make sure only one Conn is used at a time
	ts.lock.Lock()
	defer ts.lock.Unlock()
	if ts.Conn == nil {
		return fmt.Errorf("trading websocket is not connected")
	}
	err := ts.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Send Trading message: %s", string(message))
	return nil
}

func (ts *Service) reconnectSocket() {
	ts.lock.Lock()
	if ts.IsReconnecting {
		//is reconnect -> no action ->unclock -> done
		ts.lock.Unlock()
		return
	}
	ts.IsReconnecting = true
	conn := ts.Conn
	ts.Conn = nil
	ts.lock.Unlock()

	//Reset all params
	defer func() {
		ts.lock.Lock()
		ts.IsReconnecting = false
		ts.lock.Unlock()
	}()

	//Close connect if exists first
	if conn != nil {
		_ = conn.Close()
	}

	log.Println("Reconnecting to Trading Socket server after disconnect...")
	time.Sleep(5 * time.Second)
	go ts.ConnectWebsocket()
}
