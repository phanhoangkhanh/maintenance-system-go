package broadcastsocket

import (
	"maintenance-system-go/config"
	"maintenance-system-go/tradingSocket"
)
type BroadCastSocket struct {
	Service    *Service
	
}


func InitBroadCast(config *config.Config,  BridgeChannel chan tradingSocket.MarketPriceInfo) *BroadCastSocket {

	service := &Service{
		Config:   config,
		BridgeChannel: BridgeChannel,
	}

	return &BroadCastSocket{
		Service:    service,
	}
}