package models

import (
	"encoding/json"
	"github.com/antoniodipinto/ikisocket"
	"strconv"
)

type SocketEvent struct {
	Type   string      `json:"type"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

var SocketInstance *ikisocket.Websocket
var SocketUsers = make(map[string]string, 0)

func Emit(socketEvent SocketEvent, id uint) error {
	socketUserID := strconv.FormatUint(uint64(id), 10)

	if uuid, found := SocketUsers[socketUserID]; found {
		event, err := json.Marshal(socketEvent)
		if err != nil {
			return err
		}

		emitSocketErr := SocketInstance.EmitTo(uuid, event)
		if emitSocketErr != nil {
			return err
		}
	}
	return nil
}
