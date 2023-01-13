package bot

import (
	ddp "ddp"
	"reflect"
)

type statusListener func(StatusType)

// OnStatus in order to receive status changes.
func (c *Client) OnStatus(listener statusListener) {
	if c.events["OnStatus"] == nil {
		c.ddpClient.AddStatusListener(c.eventDDPHandler)
	}

	c.events["OnStatus"] = append(c.events["OnStatus"], listener)
}

// Will be emited by ddpClient
func (e *eventDDPHandler) Status(status int) {
	var st StatusType
	switch status {
	case ddp.DIALING:
		st = CONNECTING
	case ddp.CONNECTING:
		st = CONNECTING
	case ddp.CONNECTED:
		st = CONNECTED
	case ddp.DISCONNECTING:
		st = DISCONNECTING
	case ddp.DISCONNECTED:
		st = DISCONNECTED
	case ddp.RECONNECTING:
		st = RECONNECTING
	default:
		st = DISCONNECTED
	}

	e.client.updateStatus(st)
}

// informs all status listeners with the new client status.
func (c *Client) updateStatus(status StatusType) {
	if c.Status == status {
		return
	}
	c.Status = status
	for _, event := range c.events["OnStatus"] {

		// call event(status)
		f := reflect.TypeOf(event)
		if f.Kind() == reflect.Func { //is function
			if f.NumIn() == 1 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
				if f.In(0).Kind() == reflect.String { //parameter 0 is of type string (StatusType)
					go reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(status)})
				}
			}
		}
	}
}
