package bot

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"
	"github.com/ITLab-CC/bigbluebutton-bot/pad"

	convert "github.com/benpate/convert"
)

func getCookieByName(cookies []*http.Cookie, name string) string {
	result := ""
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return result
}

func (c *Client) CreateCapture(language string) (*pad.Pad, error) {
	//Subscribe to captions, pads and pads-sessions
	//Subscribe to captions
	if err := c.ddpSubscribe(bbb.CaptionsSub, nil); err != nil {
		return nil, err
	}
	captionsCollection := c.ddpClient.CollectionByName("captions")
	captionsCollection.AddUpdateListener(c.ddpEventHandler)

	//Subscribe to pads
	if err := c.ddpSubscribe(bbb.PadsSub, nil); err != nil {
		return nil, err
	}
	padsCollection := c.ddpClient.CollectionByName("pads")
	padsCollection.AddUpdateListener(c.ddpEventHandler)

	//Subscribe to pads-sessions
	if err := c.ddpSubscribe(bbb.PadsSessionsSub, nil); err != nil {
		return nil, err
	}
	padsSessionsCollection := c.ddpClient.CollectionByName("pads-sessions")
	padsSessionsCollection.AddUpdateListener(c.ddpEventHandler)

	//Create caption and add this bot as owner to it
	_, err := c.ddpCall(bbb.CreateGroupCall, "en", "captions", "English")
	if err != nil {
		return nil, err
	}

	_, err = c.ddpCall(bbb.UpdateCaptionsOwnerCall, "en", "English")
	if err != nil {
		return nil, err
	}

	_, err = c.ddpCall(bbb.CreateSessionCall, "en")
	if err != nil {
		return nil, err
	}

	//Get padID
	var padId string
	getPadIDtry := 0
	for {
		getPadIDtry++
		result, err := c.ddpCall(bbb.GetPadIdCall, "en")
		if err != nil {
			return nil, err
		}

		if getPadIDtry > 10 {
			return nil, errors.New("timeout to call getPadId: " + err.Error())
		}

		if result == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		padId = result.(string)
		break
	}
	fmt.Println("padID: " + padId)

	//Get sessionID
	var sessionID string
	getsessionIDtry := 0
	loop := true
	for loop {
		getsessionIDtry++
		result := padsSessionsCollection.FindAll()

		for _, element0 := range result {
			if element1, found := element0["sessions"]; found {
				if reflect.TypeOf(element1).Kind() == reflect.Slice {
					s := reflect.ValueOf(element1)
					for i := 0; i < s.Len(); i++ {
						element2 := s.Index(i)
						if element2.Kind() == reflect.Interface { //is Interface
							element3 := reflect.ValueOf(element2.Interface())
							if element3.Kind() == reflect.Map {
								for _, e := range element3.MapKeys() {
									element4 := element3.MapIndex(e)
									if convert.String(e) == "en" {
										sessionID = convert.String(element4)
										loop = false
									}
								}
							}
						}
					}
				}
			}
		}
		time.Sleep(100 * time.Millisecond)

		if getsessionIDtry > 100 {
			return nil, errors.New("timeout to get sessionID")
		}
	}
	fmt.Println("sessionID: " + sessionID)

	capturePad := pad.NewPad(c.PadURL, c.PadWSURL, c.SessionToken, padId, sessionID, c.SessionCookie)
	if err := capturePad.Connect(); err != nil {
		return nil, err
	}
	return capturePad, nil
}

type captureListener func()

// OnCapture in order to receive Capture changes.
func (c *Client) OnCapture(language string, listener captureListener) error {
	if c.events["OnCapture"] == nil {
		c.CreateCapture(language)
	}

	c.events["OnCapture"] = append(c.events["OnCapture"], listener)

	return nil
}

// informs all listeners with the new infos.
func (c *Client) updateCapture() {
	// // Inform all listeners
	// for _, event := range c.events["OnCapture"] {

	// 	// call event(infos)
	// 	f := reflect.TypeOf(event)
	// 	if f.Kind() == reflect.Func { //is function
	// 		if f.NumIn() == 1 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
	// 			if f.In(0).Kind() == reflect.Struct { //parameter 0 is of type string (string){ //parameter 3 is of type struct (ddp.Update)
	// 				go reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(msg)})
	// 			}
	// 		}
	// 	}
	// }
}
