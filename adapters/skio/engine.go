package skio

import (
	"fmt"
	"github.com/SeaCloudHub/notification-hub/pkg/mycontext"
	"sync"

	"github.com/SeaCloudHub/notification-hub/domain/identity"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/labstack/echo/v4"
)

type RealtimeEngine interface {
	UserSockets(userId string) []AppSocket
	EmitToRoom(room string, key string, data interface{}) error
	EmitToUser(userId string, key string, data interface{}) error
}

type rtEngine struct {
	server  *socketio.Server
	storage map[string][]AppSocket
	locker  *sync.RWMutex
}

func NewEngine() *rtEngine {
	return &rtEngine{
		storage: make(map[string][]AppSocket),
		locker:  new(sync.RWMutex),
	}
}

func (engine *rtEngine) saveAppSocket(userId string, appSck AppSocket) {
	engine.locker.Lock()

	if v, ok := engine.storage[userId]; ok {
		engine.storage[userId] = append(v, appSck)
	} else {
		engine.storage[userId] = []AppSocket{appSck}
	}

	engine.locker.Unlock()
}

func (engine *rtEngine) getAppSocket(userId string) []AppSocket {
	engine.locker.RLock()

	defer engine.locker.RUnlock()

	return engine.storage[userId]
}

func (engine *rtEngine) removeAppSocket(userId string, appSck AppSocket) {
	engine.locker.Lock()
	defer engine.locker.Unlock()

	if v, ok := engine.storage[userId]; ok {
		for i := range v {
			if v[i] == appSck {
				engine.storage[userId] = append(v[:i], v[i+1:]...)
				break
			}
		}
	}
}

func (engine *rtEngine) UserSockets(userId string) []AppSocket {
	var sockets []AppSocket

	if scks, ok := engine.storage[userId]; ok {
		return scks
	}

	return sockets
}

func (engine *rtEngine) EmitToRoom(room string, key string, data interface{}) error {
	engine.server.BroadcastToRoom("/", room, key, data)
	return nil
}

func (engine *rtEngine) EmitToUser(userId string, key string, data interface{}) error {
	sockets := engine.getAppSocket(userId)

	for _, s := range sockets {
		s.Emit(key, data)
	}

	return nil
}

func (engine *rtEngine) RunWithEchoCtx(eCtx echo.Context, identitySvc identity.Service) (error, *socketio.Server) {
	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{websocket.Default},
	})

	if err != nil {
		return err, server
	}

	engine.server = server

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Print("connecting")
		s.SetContext("")
		return nil
	})

	server.OnError("/", func(c socketio.Conn, err error) {
		fmt.Errorf("rtEngine.Run.OnError: ", err)
	})

	server.OnDisconnect("/", func(c socketio.Conn, s string) {
		fmt.Errorf("disconnected")
	})

	server.OnEvent("/", "authenticate", func(s socketio.Conn, token string) {

		iden, err := identitySvc.WhoAmI(mycontext.BuildEchoContextWithToken(eCtx, token), token)
		if err != nil {
			s.Emit("authentication_failure", err)
			s.Close()
			return
		}

		appSck := NewAppSocket(s)

		engine.saveAppSocket(iden, appSck)

		s.Emit("authenticated", iden)
	})

	go server.Serve()

	return nil, server

}

func (engine *rtEngine) Run(e *echo.Echo, identitySvc identity.Service) error {

	e.GET("/socket.io/", func(c echo.Context) error {
		_, server := engine.RunWithEchoCtx(c, identitySvc)
		server.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	return nil

}
