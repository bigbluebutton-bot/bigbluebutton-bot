package pad

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	convert "github.com/benpate/convert"
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
	Text     string
	Attribs  string
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

		Client: nil,

		AuthorID: "",
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
	if err := p.Client.On("message", p.onInitMessage); err != nil {
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

// Disconnect from the pad
func (p *Pad) Disconnect() {
	p.Client.Close()
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

func (p *Pad) onInitMessage(h *goSocketio.Channel, args ReceveClientReady) {
	if p.AuthorID == "" {
		p.AuthorID = args.Data.UserID
		p.Text = args.Data.CollabClientVars.InitialAttributedText.Text
		p.Attribs = args.Data.CollabClientVars.InitialAttributedText.Attribs
		fmt.Println("author:", p.AuthorID)
		fmt.Println("text:", p.Text)
		fmt.Println("attribs:", p.Attribs)

		//Override onInitMessage with onMessage
		p.Client.On("message", p.onMessage)
	}
}

func (p *Pad) onMessage(h *goSocketio.Channel, mapData interface{}) {
	fmt.Println(mapData)
	// Convert map to json string
	jsonStr, err := json.Marshal(mapData)
	if err != nil {
		return
	}

	// Disconnect if error/disconnect/accessStatus deny is returned
	if strings.Contains(string(jsonStr), "disconnect") || strings.Contains(string(jsonStr), "accessStatus") {
		p.Disconnect()
		return
	}

	// Convert json string to struct
	var datatype ReceveData
	if err := json.Unmarshal(jsonStr, &datatype); err != nil {
		return
	}

	// Switch datatype
	fmt.Println(datatype.Data.Type)

}

/*
A changeset describes the difference between two revisions of a document1. It has a format like this:

Z:oldLen>diffLen*ops+charBank$

where oldLen and diffLen are the lengths of the old and the difference between the old and new texts
if the diference between oldLen and newLen is >= 0, then the operation betwenn oldLen and diffLen ist >
else the operation betwenn oldLen and diffLen ist <
diffLen will always be positive!

ops are a series of operations to transform
the old text into the new text, and charBank is a string of characters that are inserted by the operations1.

The operations can be one of these types:

=: keep a number of characters unchanged
-: delete a number of characters
+: insert a number of characters from the charBank
*: apply an attribute change to a number of characters
Each operation has an optional parameter that specifies how many characters it affects. If omitted, it
defaults to 11.

For example, your changeset:

Z:1>1*0+1$h

means that you start with a text of length 1,
end with a text of length 1,
apply an attribute change (*0) to 1 character (the default),
and insert 1 character (+1) from the charBank (h).
The result is that you replace whatever character was there before with an h with some attribute.
*/
func generateChangeset(oldText *string, newText string) string {
	oldLen := len(*oldText)
	newTextLen := len(newText)
	opCode := ""
	newLen := oldLen + newTextLen

	if oldLen != 0 {
		if (*oldText)[oldLen-1] == '\n' {
			if oldLen-1 > 0 {
				opCode += "=" + convert.String(oldLen-1)
			}
			newLen -= 1
		} else {
			opCode += "=" + convert.String(oldLen)
		}
	}
	opCode += "*0"

	opCode += "+" + convert.String(newTextLen)


	opDiff := ">"

	diff := newTextLen

	*oldText = strings.Replace(*oldText, "\n", "", -1) + newText + "\n"

	result := "Z:" +
		convert.String(oldLen) + //number: old legth
		convert.String(opDiff) + //char: > or <
		convert.String(diff) + //number: diffence between old and new length
		convert.String(opCode) + //string: =, -, +, * and number
		"$" +
		convert.String(newText) //string: new text

	return result
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

	changeset := generateChangeset(&p.Text, text)
	fmt.Println(changeset)

	commandTyping := padTyping{
		Type:      "COLLABROOM",
		Component: "pad",
		Data: padTypingData{
			Type:      "USER_CHANGES",
			BaseRev:   0,
			Changeset: changeset, //"Z:1>1*0+1$g",
			Apool: padTypingDataApool{
				NumToAttrib: map[string][]string{
					"0": {"author", p.AuthorID},
				},
				NextNum: 1,
			},
		},
	}
	p.Client.Emit("message", commandTyping)


	return nil
}
