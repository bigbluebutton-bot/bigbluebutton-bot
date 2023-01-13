package bot

import (
	ddp "ddp"
	"errors"
	"reflect"
	"time"

	bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"
	"github.com/benpate/convert"
)

type groupChatMsgListener func(msg bbb.Message)

// OnGroupChatMsg in order to receive GroupChatMsg changes.
func (c *Client) OnGroupChatMsg(listener groupChatMsgListener) error {
	if c.events["OnGroupChatMsg"] == nil {

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

		collection.AddUpdateListener(c.eventDDPHandler)
	}

	c.events["OnGroupChatMsg"] = append(c.events["OnGroupChatMsg"], listener)

	return nil
}

// informs all listeners with the new infos.
func (c *Client) updateGroupChatMsg(collection string, operation string, id string, doc ddp.Update) {
	if doc == nil || doc["id"] == nil {
		return
	}
	msg := bbb.ConvertInToMessage(doc)

	// Inform all listeners
	for _, event := range c.events["OnGroupChatMsg"] {

		// call event(infos)
		f := reflect.TypeOf(event)
		if f.Kind() == reflect.Func { //is function
			if f.NumIn() == 1 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
				if f.In(0).Kind() == reflect.Struct { //parameter 0 is of type string (string){ //parameter 3 is of type struct (ddp.Update)
					go reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(msg)})
				}
			}
		}
	}
}

func (c *Client) SendChatMsg(message string, chatId string) error {
	now := time.Now()
	timestemp := convert.String(now.UnixNano())

	messageSend := bbb.MessageSend{
		ID: c.UserID + timestemp[:len(timestemp)-(len(timestemp)-13)],
		Sender: bbb.MessageSendSender{
			ID:   c.UserID,
			Name: "",
			Role: "",
		},
		ChatEmphasizedText: true,
		Message:            message,
	}

	_, err := c.ddpClient.Call("sendGroupChatMsg", chatId, messageSend)
	if err != nil {
		return errors.New("could not send message: " + err.Error())
	}

	return nil
}