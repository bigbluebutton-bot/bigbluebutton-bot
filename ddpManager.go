package bot

import (
	"errors"
	"fmt"

	ddp "github.com/gopackage/ddp"

	bbb "github.com/bigbluebutton-bot/bigbluebutton-bot/bbb"
)

//--------------------------------------------------
// ddp connection
//--------------------------------------------------

func (c *Client) ddpConnect() error {
	if c.ddpClient != nil {
		return c.ddpClient.Connect()
	}
	return errors.New("ddpClient is nil")
}

func (c *Client) ddpDisconnect() {
	if c.ddpClient != nil {
		c.ddpClient.Close()
	}
}

//--------------------------------------------------
// Call a ddp method
//--------------------------------------------------

func (c *Client) ddpCall(method bbb.CallType, params ...interface{}) (interface{}, error) {
	if c.ddpClient != nil {
		callname := bbb.GetCall(method)
		result, err := c.ddpClient.Call(callname, params...)
		if err != nil {
			return result, errors.New("could not call " + callname + ": " + err.Error())
		}
		return result, nil
	}
	return nil, errors.New("ddpClient is nil")
}

//--------------------------------------------------
// EVENTS (Collections)
//--------------------------------------------------

// This is for all events that are in "event_....go" files
type updaterfunc func(collection string, operation string, id string, doc ddp.Update)
type ddpEventHandler struct {
	client  *Client
	updater map[string][]updaterfunc
}

// Will be emited by ddpClient
func (e *ddpEventHandler) CollectionUpdate(collection string, operation string, id string, doc ddp.Update) {
	fmt.Print("CollectionUpdate: " + collection + " " + operation + " " + id + " ")
	fmt.Println(doc)
	// "redirect" to the event handler
	if flist, found := e.updater[collection]; found {
		for _, f := range flist {
			if f != nil {
				f(collection, operation, id, doc)
			}
		}
	}
}

// Subscribe to a ddp collection
func (c *Client) ddpSubscribe(collectionName bbb.SubType, callbackUpdater updaterfunc) error {
	// subscribe to bbb.collectionName
	subname, args := bbb.GetSub(collectionName) // get sub name and args
	err := c.ddpClient.Sub(subname, args...)
	if err != nil {
		return errors.New("could not subscribe to " + subname + ": " + err.Error())
	}
	collection := c.ddpClient.CollectionByName(subname)                                              // get the ddp collection
	collection.AddUpdateListener(c.ddpEventHandler)                                                  // add the update listener of the ddp collection
	c.ddpEventHandler.updater[subname] = append(c.ddpEventHandler.updater[subname], callbackUpdater) // add the update handler
	return nil
}

// Unsubscribe from a ddp collection
// ToDo: send ddp unsub and remove the update listener and the update handler
