package bot

import (
	ddp "ddp"
	"errors"
	"reflect"
)

type groupChatMsgListener func(collection string, operation string, id string, doc ddp.Update)

// OnGroupChatMsg in order to receive GroupChatMsg changes.
func (c *Client) OnGroupChatMsg(listener groupChatMsgListener) error {
	if c.events["OnGroupChatMsg"] == nil {
		e := &event{
			client: c,
		}

		err := c.ddpClient.Sub("group-chat")
		if err != nil {
			return errors.New("could not subscribe to group-chat: " + err.Error())
		}

		// subscribe to "group-chat-msg"
		err = c.ddpClient.Sub("group-chat-msg", 0)
		if err != nil {
			return errors.New("could not subscribe to group-chat-msg: " + err.Error())
		}

		collection := c.ddpClient.CollectionByName("group-chat-msg")

		collection.AddUpdateListener(e)
	}

	c.events["OnGroupChatMsg"] = append(c.events["OnGroupChatMsg"], listener)

	return nil
}

// Will be emited by ddpClient
func (e *event) CollectionUpdate(collection string, operation string, id string, doc ddp.Update) {
	e.client.updateGroupChatMsg(collection, operation, id, doc)
}

// informs all listeners with the new infos.
func (c *Client) updateGroupChatMsg(collection string, operation string, id string, doc ddp.Update) {
	// Inform all listeners
	for _, event := range c.events["OnGroupChatMsg"] {

		// call event(infos)
		f := reflect.TypeOf(event)
		if f.Kind() == reflect.Func { //is function
			if f.NumIn() == 4 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
				if (f.In(0).Kind() == reflect.String && //parameter 0 is of type string (string)
				   f.In(1).Kind() == reflect.String && //parameter 1 is of type string (string)
				   f.In(2).Kind() == reflect.String && //parameter 2 is of type string (string)
				   f.In(3).Kind() == reflect.Map) { //parameter 3 is of type struct (ddp.Update)
					reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(collection), reflect.ValueOf(operation), reflect.ValueOf(id), reflect.ValueOf(doc)})
				}
			}
		}
	}
}
