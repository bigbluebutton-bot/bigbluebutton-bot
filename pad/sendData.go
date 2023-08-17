package pad

// Send this to the Server to send a char
// {"type":"COLLABROOM","component":"pad","data":{"type":"USER_CHANGES","baseRev":0,"changeset":"Z:1>1*0+1$h","apool":{"numToAttrib":{"0":["author","a.3JMUunbWzLnaV1Ox"]},"nextNum":1}}}
type SendChar struct {
	Type      string `json:"type"`
	Component string `json:"component"`
	Data      struct {
		Type      string `json:"type"`
		BaseRev   int    `json:"baseRev"`
		Changeset string `json:"changeset"`
		Apool     struct {
			NumToAttrib struct {
				Num0 []string `json:"0"`
			} `json:"numToAttrib"`
			NextNum int `json:"nextNum"`
		} `json:"apool"`
	} `json:"data"`
}

// Send this to the Server to send a cursor position
// {"type":"COLLABROOM","component":"pad","data":{"type":"cursor","action":"cursorPosition","locationY":0,"locationX":1,"padId":"g.VPluJJUveQlgElgN$notes","myAuthorId":"a.3JMUunbWzLnaV1Ox"}}
type SendCursorPosition struct {
	Type      string `json:"type"`
	Component string `json:"component"`
	Data      struct {
		Type       string `json:"type"`
		Action     string `json:"action"`
		LocationY  int    `json:"locationY"`
		LocationX  int    `json:"locationX"`
		PadID      string `json:"padId"`
		MyAuthorID string `json:"myAuthorId"`
	} `json:"data"`
}