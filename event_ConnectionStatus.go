package bot

import ddp "github.com/gopackage/ddp"

type statusListener func(Status)

// AddStatusListener in order to receive status change updates.
func (c *Client) OnStatus(listener statusListener) {
	c.eventsOnStatus = append(c.eventsOnStatus, listener)
}

// status updates all status listeners with the new client status.
func (c *Client) Status(status int) {

	var st Status
	switch status {
		case ddp.CONNECTING:
			st = CONNECTING
		case ddp.CONNECTED:
			st = CONNECTED
		case ddp.DISCONNECTED:
			st = DISCONNECTED
		default:
			st = DISCONNECTED
	}

	if c.connectionStatus == st {
		return
	}
	c.connectionStatus = st
	for _, event := range c.eventsOnStatus {
		go event(st)
	}
}