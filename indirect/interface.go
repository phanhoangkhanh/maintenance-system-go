package indirect

import "maintenance-system-go/app"



type AppContainer interface {
	ReturnItself() *app.App
}


