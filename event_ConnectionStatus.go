package bot

import ddp "ddp"

type event struct {
	client *Client
	// statusListeners will be informed when the connection status of the client changes
	eventsOnStatus []statusListener
}

func (e *event) Status(status int) {
	var st StatusType
	switch status {
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

type statusListener func(StatusType)

// AddStatusListener in order to receive status change updates.
func (c *Client) OnStatus(listener statusListener) {
	c.event.eventsOnStatus = append(c.event.eventsOnStatus, listener)
}

// status updates all status listeners with the new client status.
func (c *Client) updateStatus(status StatusType) {
	if c.Status == status {
		return
	}
	c.Status = status
	for _, event := range c.event.eventsOnStatus {
		go event(status)
	}
}