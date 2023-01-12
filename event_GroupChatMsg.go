package bot

import (
	ddp "ddp"
	"errors"
	"reflect"

	bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"
)

type groupChatMsgListener func(msg bbb.Message)

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

	if doc == nil || doc["id"] == nil {
		return
	}
	msg := bbb.ConvertInToMessage(doc)

	e.client.updateGroupChatMsg(msg)
}

// informs all listeners with the new infos.
func (c *Client) updateGroupChatMsg(msg bbb.Message) {
	// Inform all listeners
	for _, event := range c.events["OnGroupChatMsg"] {

		// call event(infos)
		f := reflect.TypeOf(event)
		if f.Kind() == reflect.Func { //is function
			if f.NumIn() == 1 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
				if f.In(0).Kind() == reflect.Struct { //parameter 0 is of type string (string){ //parameter 3 is of type struct (ddp.Update)
					reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(msg)})
				}
			}
		}
	}
}
