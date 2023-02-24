package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"time"

	bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"

	convert "github.com/benpate/convert"
	goSocketio "github.com/graarh/golang-socketio"
	goSocketioTransport "github.com/graarh/golang-socketio/transport"
	"golang.org/x/net/publicsuffix"
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

func (c *Client) CreateCapture(language string) error {
	//Subscribe to captions, pads and pads-sessions
	//Subscribe to captions
	if err := c.ddpSubscribe(bbb.CaptionsSub, nil); err != nil {
		return err
	}
	captionsCollection := c.ddpClient.CollectionByName("captions")
	captionsCollection.AddUpdateListener(c.ddpEventHandler)

	//Subscribe to pads
	if err := c.ddpSubscribe(bbb.PadsSub, nil); err != nil {
		return err
	}
	padsCollection := c.ddpClient.CollectionByName("pads")
	padsCollection.AddUpdateListener(c.ddpEventHandler)

	//Subscribe to pads-sessions
	if err := c.ddpSubscribe(bbb.PadsSessionsSub, nil); err != nil {
		return err
	}
	padsSessionsCollection := c.ddpClient.CollectionByName("pads-sessions")
	padsSessionsCollection.AddUpdateListener(c.ddpEventHandler)

	//Create caption and add this bot as owner to it
	_, err := c.ddpCall(bbb.CreateGroupCall, "en", "captions", "English")
	if err != nil {
		return err
	}

	_, err = c.ddpCall(bbb.UpdateCaptionsOwnerCall, "en", "English")
	if err != nil {
		return err
	}

	_, err = c.ddpCall(bbb.CreateSessionCall, "en")
	if err != nil {
		return err
	}

	//Get padID
	var padId string
	getPadIDtry := 0
	for {
		getPadIDtry++
		result, err := c.ddpCall(bbb.GetPadIdCall, "en")
		if err != nil {
			return err
		}

		if getPadIDtry > 10 {
			return errors.New("timeout to call getPadId: " + err.Error())
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
			return errors.New("timeout to get sessionID")
		}
	}
	fmt.Println("sessionID: " + sessionID)

	httpclient := new(http.Client)
	//"https://example.com/pad/auth_session?padName="+padId+"&sessionID="+sessionID+"&lang=en&rtl=false&sessionToken="+c.SessionToken
	req, _ := http.NewRequest("GET", c.PadURL+"auth_session?padName="+padId+"&sessionID="+sessionID+"&lang=en&rtl=false&sessionToken="+c.SessionToken, nil)
	for _, cookie := range c.SessionCookie {
		req.AddCookie(cookie)
	}
	resp, err := httpclient.Do(req) //send request
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("pad auth_session: Server returned: " + resp.Status)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.SessionCookie = append(c.SessionCookie, &http.Cookie{Name: "sessionID", Value: sessionID}) //add sessionID cookies

	// time.Sleep(2 * time.Second)

	// type padCurserLocationData struct {
	// 	Type       string `json:"type"`
	// 	Action     string `json:"action"`
	// 	LocationY  int    `json:"locationY"`
	// 	LocationX  int    `json:"locationX"`
	// 	PadID      string `json:"padId"`
	// 	MyAuthorID string `json:"myAuthorId"`
	// }
	// type padCurserLocation struct {
	// 	Type      string                `json:"type"`
	// 	Component string                `json:"component"`
	// 	Data      padCurserLocationData `json:"data"`
	// }

	// commandCurser := `{"type":"COLLABROOM","component":"pad","data":{"type":"cursor","action":"cursorPosition","locationY":0,"locationX":1,"padId":"g.m5RsJO1Z7Rl7shlu$en","myAuthorId":"a.jV6yISPJv9ZOa9kf"}}`
	// commandCurser := padCurserLocation {
	// 	Type: "COLLABROOM",
	// 	Component: "pad",
	// 	Data: padCurserLocationData {
	// 		Type: "cursor",
	// 		Action: "cursorPosition",
	// 		LocationY: 0,
	// 		LocationX: 1,
	// 		PadID: padId,
	// 		MyAuthorID: "a.jV6yISPJv9ZOa9kf",
	// 	},
	// }
	// commandCurser := `{
	// 	"type":"COLLABROOM",
	// 	"component":"pad",
	// 	"data":{
	// 	   "type":"USER_CHANGES",
	// 	   "baseRev":0,
	// 	   "changeset":"Z:1>1*0+1$h",
	// 	   "apool":{
	// 		  "numToAttrib":{
	// 			 "0":[
	// 				"author",
	// 				"a.jV6yISPJv9ZOa9kf"
	// 			 ]
	// 		  },
	// 		  "nextNum":1
	// 	   }
	// 	}
	//  }`

	uri := c.PadWSURL + "socket.io/?sessionToken=" + c.SessionToken + "&padId=" + padId + "&EIO=3&transport=websocket"

	transp := goSocketioTransport.GetDefaultWebsocketTransport()

	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	ur, _ := url.Parse(c.PadURL)
	jar.SetCookies(ur, c.SessionCookie)

	transp.Cookie = jar

	client := goSocketio.NewClient()

	err = client.On(goSocketio.OnDisconnection, func(h *goSocketio.Channel) {
		fmt.Println("Disconnected")
	})
	if err != nil {
		return err
	}

	err = client.On(goSocketio.OnConnection, func(h *goSocketio.Channel) {
		fmt.Println("Connected")

		time.Sleep(2 * time.Second)
		//212:42["message",{"component":"pad","type":"CLIENT_READY","padId":"g.9d4O2LRqTkIfh6bM$notes","sessionID":"s.4918c0b0b9b7913b5e29334a50f58212","token":"t.oNTJCeHhA5x2lI9rM5st","userInfo":{"colorId":null,"name":null}}]
		token := "t." + c.SessionToken //token can be anything. So we take the sessionToken
		jsonStr := `{"component":"pad","type":"CLIENT_READY","padId":"` + padId + `","sessionID":"` + getCookieByName(c.SessionCookie, "sessionID") + `","token":"` + token + `","userInfo":{"colorId":null,"name":null}}`
		type ClientReady struct {
			Component string      `json:"component"`
			Type      string      `json:"type"`
			PadID     string      `json:"padId"`
			SessionID string      `json:"sessionID"`
			Token     string      `json:"token"`
			UserInfo  interface{} `json:"userInfo"`
		}
		var cr ClientReady
		err = json.Unmarshal([]byte(jsonStr), &cr)
		if err != nil {
			panic(err)
		}
		client.Emit("message", cr)
	})
	if err != nil {
		return err
	}

	type clientReadyResponse struct {
		Data struct {
			UserID string `json:"userId"`
		} `json:"data"`
	}

	err = client.On("message", func(h *goSocketio.Channel, args clientReadyResponse) {
		author := args.Data.UserID
		fmt.Println("author:", author)
		time.Sleep(2 * time.Second)

		type padTypingDataApool struct {
			NumToAttrib map[string][]string `json:"numToAttrib"`
			NextNum     int                 `json:"nextNum"`
		}
		type padTypingData struct {
			Type      string             `json:"type"`
			BaseRev   int                `json:"baseRev"`
			Changeset string             `json:"changeset"`
			Apool     padTypingDataApool `json:"apool"`
		}
		type padTyping struct {
			Type      string        `json:"type"`
			Component string        `json:"component"`
			Data      padTypingData `json:"data"`
		}
		commandTyping := padTyping{
			Type:      "COLLABROOM",
			Component: "pad",
			Data: padTypingData{
				Type:      "USER_CHANGES",
				BaseRev:   0,
				Changeset: "Z:1>1*0+1$h",
				Apool: padTypingDataApool{
					NumToAttrib: map[string][]string{
						"0": []string{"author", author},
					},
					NextNum: 1,
				},
			},
		}
		client.Emit("message", commandTyping)

	})
	if err != nil {
		return err
	}

	err = client.Dial(
		uri,
		transp)
	if err != nil {
		return err
	} else {
		fmt.Println("Connecting...")
	}

	//https://example.com/pad/api/1.2.14/appendText?apikey=JxkIWG6YFOmbT3xEmhDaG42K9cXUkqs0vh6BGPHi8ksMP3VjsTgvC9H2yTRkW&padID=g.jUdjji2zr14keg5Y$en&text=oh%20no

	return nil
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
