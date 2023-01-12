package bot

import (
	ddp "ddp"
	"reflect"
)

type templateListener func(string)

// OnTemplate in order to receive Template changes.
func (c *Client) OnTemplate(listener templateListener) {
	if c.events["OnTemplate"] == nil {
		e := &event{
			client: c,
		}
		c.ddpClient.AddTemplateListener(e)
	}

	c.events["OnTemplate"] = append(c.events["OnStatus"], listener)
}

// Will be emited by ddpClient
func (e *event) Template(infos string) {
	// Do stuff

	e.client.updateStatus(infos)
}

// informs all listeners with the new infos.
func (c *Client) updateTemplate(infos string) {
	// Inform all listeners
	for _, event := range c.events["OnTemplate"] {

		// call event(infos)
		f := reflect.TypeOf(event)
		if f.Kind() == reflect.Func { //is function
			if f.NumIn() == 1 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
				if f.In(0).Kind() == reflect.String { //parameter 0 is of type string (StatusType)
					reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(status)})
				}
			}
		}
	}
}
