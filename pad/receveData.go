package pad

// Initial response from the server
type ReceveClientReady struct {
	Data struct {
		UserID string `json:"userId"`
		CollabClientVars struct {
			InitialAttributedText struct {
				Text    string `json:"text"`
				Attribs string `json:"attribs"`
			} `json:"initialAttributedText"`
			//...
		} `json:"collab_client_vars"`
		//...
	} `json:"data"`
}


// Server will send data which has a Type. To get the type use ReceveData.Data.Type
// {"type":"COLLABROOM","data":{"type":"ACCEPT_COMMIT",...}}
type ReceveData struct {
	Type string `json:"type"`
	Data struct {
		Type   string `json:"type"`
	} `json:"data"`
}

// Server will confim the SendChar sent by this client
// {"type":"COLLABROOM","data":{"type":"ACCEPT_COMMIT","newRev":1}}
type ReceveConfirmSendChar struct {
	Type string `json:"type"`
	Data struct {
		Type   string `json:"type"`
		NewRev int    `json:"newRev"`
	} `json:"data"`
}

// Server will send char, if a outher user wrote a char
//{"type":"COLLABROOM","data":{"type":"NEW_CHANGES","newRev":3,"changeset":"Z:3>1=2*0+1$b","apool":{"numToAttrib":{"0":["author","a.MO7GXKUWttjc4se8"]},"attribToNum":{"author,a.MO7GXKUWttjc4se8":0},"nextNum":1},"author":"a.MO7GXKUWttjc4se8","currentTime":1677492927116,"timeDelta":null}}
type ReceveSendChar struct {
	Type string `json:"type"`
	Data struct {
		Type      string `json:"type"`
		NewRev    int    `json:"newRev"`
		Changeset string `json:"changeset"`
		Apool     struct {
			NumToAttrib struct {
				Num0 []string `json:"0"`
			} `json:"numToAttrib"`
			AttribToNum interface{} `json:"attribToNum"`
			NextNum int `json:"nextNum"`
		} `json:"apool"`
		Author      string `json:"author"`
		CurrentTime int64  `json:"currentTime"`
		TimeDelta   any    `json:"timeDelta"`
	} `json:"data"`
}

// Server will send a cursor position, if a outher user moved the cursor
// {"type":"COLLABROOM","data":{"type":"CUSTOM","payload":{"action":"cursorPosition","authorId":"a.3JMUunbWzLnaV1Ox","authorName":"Julian","padId":"g.VPluJJUveQlgElgN$notes","locationX":0,"locationY":0}}}
type ReceveCursorPosition struct {
	Type string `json:"type"`
	Data struct {
		Type    string `json:"type"`
		Payload struct {
			Action     string `json:"action"`
			AuthorID   string `json:"authorId"`
			AuthorName string `json:"authorName"`
			PadID      string `json:"padId"`
			LocationX  int    `json:"locationX"`
			LocationY  int    `json:"locationY"`
		} `json:"payload"`
	} `json:"data"`
}

// Server will send a new user, if a outher user joined the pad
// {"type":"COLLABROOM","data":{"type":"USER_NEWINFO","userInfo":{"colorId":9,"name":"Julian","userId":"a.MO7GXKUWttjc4se8"}}}
type ReceveNewUser struct {
	Type string `json:"type"`
	Data struct {
		Type     string `json:"type"`
		UserInfo struct {
			ColorID int    `json:"colorId"`
			Name    string `json:"name"`
			UserID  string `json:"userId"`
		} `json:"userInfo"`
	} `json:"data"`
}