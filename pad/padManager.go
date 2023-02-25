package pad

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	goSocketio "github.com/graarh/golang-socketio"
	goSocketioTransport "github.com/graarh/golang-socketio/transport"
	"golang.org/x/net/publicsuffix"
)

type Pad struct {
	URL          string //"https://example.com/pad/"
	WsURL        string //"wss://example.com/pad/"
	SessionToken string
	PadId        string
	SessionID    string
	Cookie       []*http.Cookie

	Client *goSocketio.Client

	AuthorID string
}

// Create new pad
// wsURL = "wss://example.com/pad/"
// SessionToken = gtxiomrffih2b8qr (from bbb)
// padId = g.9d4O2LRqTkIfh6bM$notes (from ddp. To get it c.ddpCall(bbb.GetPadIdCall, "en"))
// sessionID = s.4918c0b0b9b7913b5e29334a50f58212 (from ddp. To get it padsSessionsCollection.FindAll())
// cookie = client.SessionCookie
func NewPad(url string, wsURL string, sessionToken string, padId string, sessionID string, cookie []*http.Cookie) *Pad {
	// Add sessionID cookies
	if getCookieByName(cookie, "sessionID") == "" {
		cookie = append(cookie, &http.Cookie{Name: "sessionID", Value: sessionID}) //add sessionID cookies
	}
	sessionIDvalue := getCookieByName(cookie, "sessionID")
	if !strings.Contains(sessionIDvalue, sessionID) {
		for _, cookie := range cookie {
			if cookie.Name == "sessionID" {
				cookie.Value = cookie.Value + "," + sessionID
				break
			}
		}
	}

	return &Pad{
		URL:          url,
		WsURL:        wsURL + "socket.io/?sessionToken=" + sessionToken + "&padId=" + padId + "&EIO=3&transport=websocket",
		SessionToken: sessionToken,
		PadId:        padId,
		SessionID:    sessionID,
		Cookie:       cookie,

		Client: 	  nil,

		AuthorID: 	  "",
	}
}

func getCookieByName(cookies []*http.Cookie, name string) string {
	result := ""
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return result
}

// Register session
func (p *Pad) RegisterSession() error {
	httpclient := new(http.Client)
	//"https://example.com/pad/auth_session?padName="+padId+"&sessionID="+sessionID+"&lang=en&rtl=false&sessionToken="+c.SessionToken
	req, _ := http.NewRequest("GET", p.URL+"auth_session?padName="+p.PadId+"&sessionID="+p.SessionID+"&lang=en&rtl=false&sessionToken="+p.SessionToken, nil)
	for _, cookie := range p.Cookie {
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
	return nil
}

// Connect to the pad
func (p *Pad) Connect() error {
	if err := p.RegisterSession(); err != nil {
		return err
	}

	//Create cookie jar
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	url, _ := url.Parse(p.URL)
	jar.SetCookies(url, p.Cookie)

	//Get websocket transport
	transport := goSocketioTransport.GetDefaultWebsocketTransport()
	//Set cookies
	transport.Cookie = jar

	//Create client
	p.Client = goSocketio.NewClient()

	//Create events
	//On Connection
	if err := p.Client.On(goSocketio.OnConnection, p.onConnect); err != nil {
		return err
	}
	//On Disconnection
	if err := p.Client.On(goSocketio.OnDisconnection, p.onDisconnect); err != nil {
		return err
	}
	//On message
	if err := p.Client.On("message", p.onMessage); err != nil {
		return err
	}

	//Connect to the server
	err := p.Client.Dial(
		p.WsURL,
		transport)
	if err != nil {
		return err
	} else {
		fmt.Println("Connecting...")
	}
	return nil
}

func (p *Pad) onConnect(h *goSocketio.Channel) {
	fmt.Println("Connected")

	// Send ClientReady
	//212:42["message",{"component":"pad","type":"CLIENT_READY","padId":"g.9d4O2LRqTkIfh6bM$notes","sessionID":"s.4918c0b0b9b7913b5e29334a50f58212","token":"t.oNTJCeHhA5x2lI9rM5st","userInfo":{"colorId":null,"name":null}}]
	type ClientReadyUserInfo struct {
		ColorID any `json:"colorId"`
		Name    any `json:"name"`
	}
	type ClientReady struct {
		Component string              `json:"component"`
		Type      string              `json:"type"`
		PadID     string              `json:"padId"`
		SessionID string              `json:"sessionID"`
		Token     string              `json:"token"`
		UserInfo  ClientReadyUserInfo `json:"userInfo"`
	}

	cr := ClientReady{
		Component: "pad",
		Type:      "CLIENT_READY",
		PadID:     p.PadId,
		SessionID: p.SessionID,
		Token:     "t." + p.SessionToken, //token can be anything. So we take the sessionToken
		UserInfo: ClientReadyUserInfo{
			ColorID: nil,
			Name:    nil,
		},
	}
	// Send ClientReady
	p.Client.Emit("message", cr)
}

func (p *Pad) onDisconnect(h *goSocketio.Channel) {
	fmt.Println("Disconnected")
}

type clientReadyResponse struct {
	Data struct {
		UserID string `json:"userId"`
	} `json:"data"`
}

func (p *Pad) onMessage(h *goSocketio.Channel, args clientReadyResponse) {
	if p.AuthorID == "" {
		p.AuthorID = args.Data.UserID
		fmt.Println("author:", p.AuthorID)
	}
}


func (p *Pad) SendText(text string) error {
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
			Changeset: "Z:1>1*0+1$g",
			Apool: padTypingDataApool{
				NumToAttrib: map[string][]string{
					"0": {"author", p.AuthorID},
				},
				NextNum: 1,
			},
		},
	}
	return p.Client.Emit("message", commandTyping)
}