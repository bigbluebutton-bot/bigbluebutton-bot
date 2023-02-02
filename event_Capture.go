package bot

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	// "reflect"

	// bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"
	// "bufio"
	// socketio "github.com/zhouhui8915/go-socket.io-client"
	// "log"
	// "os"
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

func(c *Client) CreateCapture(language string) error {
	//Subscribe to captions, pads and pads-sessions
	//Subscribe to captions
	err := c.ddpClient.Sub("captions")
	if err != nil {
		return errors.New("could not subscribe to captions: " + err.Error())
	}
	captionsCollection := c.ddpClient.CollectionByName("captions")
	captionsCollection.AddUpdateListener(c.eventDDPHandler)

	//Subscribe to pads
	err = c.ddpClient.Sub("pads")
	if err != nil {
		return errors.New("could not subscribe to pads: " + err.Error())
	}
	padsCollection := c.ddpClient.CollectionByName("pads")
	padsCollection.AddUpdateListener(c.eventDDPHandler)

	//Subscribe to pads-sessions
	err = c.ddpClient.Sub("pads-sessions")
	if err != nil {
		return errors.New("could not subscribe to group-chat: " + err.Error())
	}
	padsSessionsCollection := c.ddpClient.CollectionByName("pads-sessions")
	padsSessionsCollection.AddUpdateListener(c.eventDDPHandler)


	//Create caption and add this bot as owner to it
	_, err = c.ddpClient.Call("createGroup", "en", "captions", "English")
	if err != nil {
		return errors.New("could not call createGroup: " + err.Error())
	}

	_, err = c.ddpClient.Call("updateCaptionsOwner", "en", "English")
	if err != nil {
		return errors.New("could not call updateCaptionsOwner: " + err.Error())
	}

	_, err = c.ddpClient.Call("createSession", "en")
	if err != nil {
		return errors.New("could not call createSession: " + err.Error())
	}


	var padId string
	getPadIDtry := 0
	for {
		getPadIDtry++
		result, err := c.ddpClient.Call("getPadId", "en")
		if err != nil {
			return errors.New("could not call getPadId: " + err.Error())
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
	fmt.Println(padId)

	time.Sleep(2 * time.Second)

	// fmt.Println(captionsCollection.FindAll())
	// fmt.Println(padsCollection.FindAll())
	// fmt.Println(padsSessionsCollection.FindAll())

	httpclient := new(http.Client)
	//"https://example.com/pad/auth_session?padName="+padId+"&sessionID=s.42e1ba9a2c46a7d587d8b5896b625080&lang=en&rtl=false&sessionToken=fqicxhxpidc3hyiz"
	hash := md5.Sum([]byte(padId))
	sessionID := hex.EncodeToString(hash[:]) // The sessionID looks like a MD5 hash, but I dont know
	req, _ := http.NewRequest("GET", c.PadURL + "auth_session?padName=" + padId + "&sessionID=s." + sessionID + "&lang=en&rtl=false&sessionToken="+c.SessionToken, nil)
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
	c.SessionCookie = append(c.SessionCookie, &http.Cookie{Name: "sessionID", Value: "s." + sessionID}) //add sessionID cookies

	time.Sleep(2 * time.Second)
	fmt.Println(padsSessionsCollection.FindAll())



















	// //Make a http get request to the BigBlueButton API
	// httpclient = new(http.Client)
	// req, _ = http.NewRequest("GET", c.PadURL + "socket.io/?sessionToken="+c.SessionToken+"&padId="+padId+"&EIO=3&transport=polling", nil)
	// for _, cookie := range c.SessionCookie {
	// 	req.AddCookie(cookie)
	// }
	// resp, err = httpclient.Do(req) //send request
	// if err != nil {
	// 	return err
	// }
	// if resp.StatusCode != 200 {
	// 	return errors.New("Server returned: " + resp.Status)
	// }

	// defer resp.Body.Close()
	// body, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// type connection struct {
	// 	Sid          string   `json:"sid"`
	// 	Upgrades     []string `json:"upgrades"`
	// 	PingInterval int      `json:"pingInterval"`
	// 	PingTimeout  int      `json:"pingTimeout"`
	// }

	// //Unmarshal xml
	// body = body[4:]
	// //fmt.Println(body)
	// var response connection
	// err = json.Unmarshal(body, &response)
	// if err != nil {
	// 	return err
	// }
	// sid := response.Sid
	// fmt.Println(sid)


	// //--------------------------------------------------------------


	// httpclient = new(http.Client)
	// req, _ = http.NewRequest("GET", c.PadURL + "socket.io/?sessionToken="+c.SessionToken+"&padId="+padId+"&EIO=3&transport=polling&sid=" + sid, nil)
	// for _, cookie := range c.SessionCookie {
	// 	req.AddCookie(cookie)
	// }
	// resp, err = httpclient.Do(req) //send request
	// if err != nil {
	// 	return err
	// }
	// if resp.StatusCode != 200 {
	// 	return errors.New("Server returned: " + resp.Status)
	// }

	// defer resp.Body.Close()
	// body, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// token := resp.Cookies()
	// c.SessionCookie = append(c.SessionCookie, &http.Cookie{Name: "token", Value: getCookieByName(token, "token")}) //add sessionID cookies
	// fmt.Println(string(body))

	// //--------------------------------------------------------------
	// //212:42["message",{"component":"pad","type":"CLIENT_READY","padId":"g.9d4O2LRqTkIfh6bM$notes","sessionID":"s.4918c0b0b9b7913b5e29334a50f58212","token":"t.oNTJCeHhA5x2lI9rM5st","userInfo":{"colorId":null,"name":null}}]
	// PostData := strings.NewReader(`244:42["message",{"component":"pad","type":"CLIENT_READY","padId":"`+padId+`","sessionID":"`+getCookieByName(c.SessionCookie, "sessionID")+`","token":"`+getCookieByName(c.SessionCookie, "token")+`","userInfo":{"colorId":null,"name":null}}]`)

	// httpclient = new(http.Client)
	// req, _ = http.NewRequest("POST", c.PadURL + "socket.io/?sessionToken="+c.SessionToken+"&padId="+padId+"&EIO=3&transport=polling&sid=" + sid, PostData)
	// for _, cookie := range c.SessionCookie {
	// 	req.AddCookie(cookie)
	// }
	// resp, err = httpclient.Do(req) //send request
	// if err != nil {
	// 	return err
	// }
	// if resp.StatusCode != 200 {
	// 	return errors.New("Server returned: " + resp.Status)
	// }

	// defer resp.Body.Close()
	// body, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(body))

	// //--------------------------------------------------------------

	// httpclient = new(http.Client)
	// req, _ = http.NewRequest("GET", c.PadURL + "socket.io/?sessionToken="+c.SessionToken+"&padId="+padId+"&EIO=3&transport=polling&sid=" + sid, nil)
	// for _, cookie := range c.SessionCookie {
	// 	req.AddCookie(cookie)
	// 	fmt.Println(cookie)
	// }
	// resp, err = httpclient.Do(req) //send request
	// if err != nil {
	// 	return err
	// }
	// // if resp.StatusCode != 200 {
	// // 	return errors.New("Server returned: " + resp.Status)
	// // }

	// defer resp.Body.Close()
	// body, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(body))











	// // uri := c.PadURL + "socket.io/?sessionToken=" + c.SessionToken + "&padId=" + padId + "&EIO=3&transport=websocket&sid=" + sid
	// uri := c.PadURL + "socket.io/"

	// opts := &socketio_client.Options{
	// 	Transport: "websocket",
	// 	// Query:     make(map[string]string),
	// 	Query: map[string]string{
	// 		"padId":        padId,
	// 		"sessionToken": c.SessionToken,
	// 		// "EIO": "3",
	// 		// "transport": "websocket",
	// 		// "sid": sid,
	// 	},
	// 	Cookies: c.SessionCookie,
	// }

	// socketclient, err := socketio_client.NewClient(uri, opts)
	// if err != nil {
	// 	fmt.Printf("NewClient error:%v\n", err)
	// 	return err
	// }

	// socketclient.On("error", func() {
	// 	fmt.Printf("on error\n")
	// })
	// socketclient.On("connection", func() {
	// 	fmt.Printf("on connect\n")

	// })
	// socketclient.On("message", func(msg string) {
	// 	fmt.Printf("on message:%v\n", msg)
	// })
	// socketclient.On("disconnection", func() {
	// 	fmt.Printf("on disconnect\n")
	// })

	// socketclient.Connect()

	// // socketclient.Emit("message", "2probe")
	// // socketclient.Emit("message", "5")
	// time.Sleep(1 * time.Second)

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

	// //42["message",{"type":"COLLABROOM","component":"pad","data":{"type":"cursor","action":"cursorPosition","locationY":0,"locationX":1,"padId":"g.m5RsJO1Z7Rl7shlu$en","myAuthorId":"a.jV6yISPJv9ZOa9kf"}}]
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
	// // commandCurser := `{
	// // 	"type":"COLLABROOM",
	// // 	"component":"pad",
	// // 	"data":{
	// // 	   "type":"USER_CHANGES",
	// // 	   "baseRev":0,
	// // 	   "changeset":"Z:1>1*0+1$h",
	// // 	   "apool":{
	// // 		  "numToAttrib":{
	// // 			 "0":[
	// // 				"author",
	// // 				"a.jV6yISPJv9ZOa9kf"
	// // 			 ]
	// // 		  },
	// // 		  "nextNum":1
	// // 	   }
	// // 	}
	// //  }`
	// err = socketclient.Emit("message", commandCurser)
	// if err != nil {
	// 	panic(err)
	// }




	// type padTypingDataApool struct {
	// 	NumToAttrib map[string][]string `json:"numToAttrib"`
	// 	NextNum     int                 `json:"nextNum"`
	// }
	// type padTypingData struct {
	// 	Type      string             `json:"type"`
	// 	BaseRev   int                `json:"baseRev"`
	// 	Changeset string             `json:"changeset"`
	// 	Apool     padTypingDataApool `json:"apool"`
	// }
	// type padTyping struct {
	// 	Type      string        `json:"type"`
	// 	Component string        `json:"component"`
	// 	Data      padTypingData `json:"data"`
	// }

	// //42["message",{"type":"COLLABROOM","component":"pad","data":{"type":"USER_CHANGES","baseRev":0,"changeset":"Z:1>1*0+1$h","apool":{"numToAttrib":{"0":["author","a.jV6yISPJv9ZOa9kf"]},"nextNum":1}}}]
	// commandTyping := padTyping {
	// 	Type: "COLLABROOM",
	// 	Component: "pad",
	// 	Data: padTypingData {
	// 		Type: "USER_CHANGES",
	// 		BaseRev: 0,
	// 		Changeset: "Z:1>1*0+1$h",
	// 		Apool: padTypingDataApool {
	// 			NumToAttrib: map[string][]string {
	// 				"0": []string{"author", "a.jV6yISPJv9ZOa9kf"},
	// 			},
	// 			NextNum: 1,
	// 		},
	// 	},
	// }
	// err = socketclient.Emit("message", commandTyping)
	// if err != nil {
	// 	panic(err)
	// }


	// // command := `{"type":"COLLABROOM","component":"pad","data":{"type":"cursor","action":"cursorPosition","locationY":0,"locationX":0,"padId":"` + padId + `","myAuthorId":"a.oj3x18fQdseyxqNG"}}`
	// // socketclient.Emit("message", command)

	// // command = `{"type":"COLLABROOM","component":"pad","data":{"type":"cursor","action":"cursorPosition","locationY":0,"locationX":1,"padId":"` + padId + `","myAuthorId":"a.oj3x18fQdseyxqNG"}}`
	// // socketclient.Emit("message", command)

	// // command = `{"type":"COLLABROOM","component":"pad","data":{"type":"USER_CHANGES","baseRev":0,"changeset":"Z:1>1*0+1$h","apool":{"numToAttrib":{"0":["author","a.oj3x18fQdseyxqNG"]},"nextNum":1}}}`
	// // socketclient.Emit("message", command)



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
