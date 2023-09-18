package bot

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gorilla/websocket"

	// "fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"fmt"
	"log"
	"net"

	logging "github.com/pion/logging"
	turn "github.com/pion/turn/v2"
	"github.com/pion/webrtc/v3"
)

type stuns struct {
	StunServers []struct {
		URL string `json:"url"`
	} `json:"stunServers"`
	TurnServers []struct {
		Username string `json:"username"`
		Password string `json:"password"`
		URL      string `json:"url"`
		TTL      int    `json:"ttl"`
	} `json:"turnServers"`
	RemoteIceCandidates []any `json:"remoteIceCandidates"`
}

// ListenToAudio joins the audio channel of the meeting and starts listening to the audio stream.
func (c *Client) ListenToAudio() error {

	// Make request to https://example.com/bigbluebutton/api/stuns?sessionToken=TOKEN
	// to get the STUN server address
	httpclient := new(http.Client)
	req, _ := http.NewRequest("GET", c.API.Url + "stuns?sessionToken="+c.SessionToken, nil)
	// Add cookies
	for _, cookie := range c.SessionCookie {
		req.AddCookie(cookie)
	}
	
	resp, err := httpclient.Do(req) //send request
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("bbb api: Couldnt get stun server address. Server returned: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal json body to stuns struct
	var stun stuns
	err = json.Unmarshal(body, &stun)
	if err != nil {
		return err
	}

	fmt.Println(stun)

	
	// Prüfen, ob TURN-Serverinformationen vorhanden sind
	if len(stun.TurnServers) == 0 {
		return errors.New("bbb api: No TURN servers provided")
	}

	// Für Einfachheit nehmen wir den ersten TURN-Server (könnte erweitert werden, um den besten zu wählen)
	turnServer := stun.TurnServers[0]


	// Make api request to get all information of this meeting (VoiceBridge, CaleeName, UserID, UserName)
	meetings, err := c.API.GetMeetings()
	 if err != nil {
		return err
	}

	// Get the meeting by the internal meeting id
	meeting := meetings[c.ExternalMeetingID]
	voiceBridge := meeting.VoiceBridge
	caleeName := "GLOBAL_AUDIO_" + strconv.FormatInt(int64(voiceBridge), 10)

	//Connect to the signalling server
	if err := connectToWebSocketSignallingServer(c.WebRTCWSURL, c.SessionToken, c.InternalMeetingID, voiceBridge, caleeName, c.InternalUserID, c.UserName, c.SessionCookie); err != nil {
		return err
	}

	// Verwenden Sie die TURN-Informationen, um eine Verbindung herzustellen
	err = connectToTurnServer(turnServer)
	if err != nil {
		return err
	}


	return nil
}

type start struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Role              string `json:"role"`
	InternalMeetingID string `json:"internalMeetingId"`
	VoiceBridge       string `json:"voiceBridge"`
	CaleeName         string `json:"caleeName"`
	UserID            string `json:"userId"`
	UserName          string `json:"userName"`
}

type startResponse struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Response  string `json:"response"`
	SdpAnswer string `json:"sdpAnswer"`
}

type subscriberAnswer struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Role        string `json:"role"`
	VoiceBridge string `json:"voiceBridge"`
	SdpOffer    string `json:"sdpOffer"`
}

type webRTCAudioSuccess struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Success string `json:"success"`
}

type pingpong struct {
	ID string `json:"id"`
}

func connectToWebSocketSignallingServer(URL, token, internalMeetingID string, voiceBridge int, caleeName, userID, userName string, cookies []*http.Cookie) error {
	u, err := url.Parse(URL + "?sessionToken=" + token)
	if err != nil {
		return fmt.Errorf("failed to parse WebRTC WebSocket URL: %v", err)
	}
	wsurl := u.String()

	// Header für zusätzliche Optionen (optional)
	header := http.Header{}

	// Read all cookies and create a http.CookieJar
	coockieJar, err := cookiejar.New(nil)
	if err != nil { 
		return fmt.Errorf("failed to create cookie jar: %v", err)
	}
	tempu := url.URL(*u)
	tempu.Scheme = "https"
	coockieJar.SetCookies(&tempu, cookies)

	// WebSocket-Verbindung herstellen
	wsdialer := websocket.DefaultDialer
	wsdialer.Jar = coockieJar
	conn, _, err := wsdialer.Dial(wsurl, header)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to WebSocket server!")

	// Nachricht 1 senden
	msg1 := start{
		ID:                "start",
		Type:              "audio",
		Role:              "recv",
		InternalMeetingID: internalMeetingID,
		VoiceBridge:       strconv.FormatInt(int64(voiceBridge), 10),
		CaleeName:         caleeName,
		UserID:            userID,
		UserName:          userName,
	}
	err = conn.WriteJSON(msg1)
	if err != nil {
		return fmt.Errorf("failed to send message 1: %v", err)
	}
	fmt.Println("Sent message 1!")

	// Antwort 1 empfangen
	var resp1 startResponse
	err = conn.ReadJSON(&resp1)
	if err != nil {
		return fmt.Errorf("failed to read response 1: %v", err)
	}
	fmt.Println("Received response 1!")
	fmt.Println(resp1)

	// Nachricht 2 senden
	// Generate spdOffer
	sdpOffer, err := generateSDPOffer()
	if err != nil {
		return fmt.Errorf("failed to generate SDP offer: %v", err)
	}
	fmt.Println(sdpOffer)

	// Hier setzen wir den `sdpOffer` Wert statisch, aber in einer realen Anwendung würden Sie diesen Wert dynamisch generieren.
	msg2 := subscriberAnswer{
		ID:          "subscriberAnswer",
		Type:        "audio",
		Role:        "recv",
		VoiceBridge: strconv.FormatInt(int64(voiceBridge), 10),
		SdpOffer:    sdpOffer,
	}
	err = conn.WriteJSON(msg2)
	if err != nil {
		return fmt.Errorf("failed to send message 2: %v", err)
	}
	fmt.Println("Sent message 2!")

	// Antwort 2 empfangen
	var resp2 webRTCAudioSuccess
	err = conn.ReadJSON(&resp2)
	if err != nil {
		return fmt.Errorf("failed to read response 2: %v", err)
	}
	fmt.Println("Received response 2!")
	fmt.Println(resp2)

	// Ping-Nachricht senden
	pingMsg := pingpong{
		ID: "ping",
	}
	err = conn.WriteJSON(pingMsg)
	if err != nil {
		return fmt.Errorf("failed to send ping message: %v", err)
	}
	fmt.Println("Sent ping message!")

	// Pong-Antwort empfangen
	var pongResp pingpong
	err = conn.ReadJSON(&pongResp)
	if err != nil {
		return fmt.Errorf("failed to read pong response: %v", err)
	}
	fmt.Println("Received pong response!")
	fmt.Println(pongResp)

	log.Println("Successfully completed the conversation with the WebSocket server!")
	return nil
}



func generateSDPOffer() (string, error) {
	// Erstellen Sie eine neue RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return "", fmt.Errorf("failed to create new peer connection: %v", err)
	}

	// Erstellen Sie einen neuen Audio-Transceiver
	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio)
	if err != nil {
		return "", fmt.Errorf("failed to add audio transceiver: %v", err)
	}

	// Erstellen Sie ein Angebot
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create SDP offer: %v", err)
	}

	// Setzen Sie das lokale Angebot
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		return "", fmt.Errorf("failed to set local description: %v", err)
	}

	// Geben Sie das SDP-Angebot zurück
	return offer.SDP, nil
}


func connectToTurnServer(serverInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	TTL      int    `json:"ttl"`
}) error {

	// Aufteilen von URL in Host und Port
	u, err := url.Parse(serverInfo.URL)
	if err != nil {
		return fmt.Errorf("failed to parse TURN server URL: %v", err)
	}



	host, _, err := net.SplitHostPort(u.Opaque)
	if err != nil {
		return fmt.Errorf("failed to parse TURN server URL: %v", err)
	}


	// port, err := strconv.Atoi(portStr)
	// if err != nil {
	// 	return fmt.Errorf("invalid port in TURN server URL: %v", err)
	// }

	turnServerAddr := u.Opaque

	// // Aufteilen von Username und Password (angenommen, sie sind im Format "user:password")
	// cred := strings.SplitN(serverInfo.Username, ":", 2)

	// TURN client won't create a local listening socket by itself.
	conn, err := net.ListenPacket("udp4", "0.0.0.0:0")
	if err != nil {
		return fmt.Errorf("failed to create local listening socket: %v", err)
	}
	defer conn.Close()

	// turnServerAddr := fmt.Sprintf("%s:%d", host, port)

	loggerFactory := logging.NewDefaultLoggerFactory()
	loggerFactory.DefaultLogLevel = logging.LogLevelTrace // Setzen Sie das LogLevel auf Trace für detaillierte Protokolle
	cfg := &turn.ClientConfig{
		STUNServerAddr: turnServerAddr,
		TURNServerAddr: turnServerAddr,
		Conn:           conn,
		Username:       serverInfo.Username,
		Password:       serverInfo.Password,
		Realm:          host, // Verwenden Sie den Host als Realm
		LoggerFactory:  logging.NewDefaultLoggerFactory(),
	}

	client, err := turn.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create TURN client: %v", err)
	}
	defer client.Close()

	// Start listening on the conn provided.
	err = client.Listen()
	if err != nil {
		return fmt.Errorf("TURN client failed to listen: %v", err)
	}

	// Allocate a relay socket on the TURN server. On success, it
	// will return a net.PacketConn which represents the remote
	// socket.
	relayConn, err := client.Allocate()
	if err != nil {
		return fmt.Errorf("TURN client failed to allocate: %v", err)
	}
	defer relayConn.Close()

	log.Printf("Connected to TURN server. Relayed address: %s", relayConn.LocalAddr().String())

	return nil
}