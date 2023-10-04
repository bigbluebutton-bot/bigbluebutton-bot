package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	logging "github.com/pion/logging"
	"github.com/pion/sdp/v3"
	turn "github.com/pion/turn/v2"
	"github.com/pion/webrtc/v3"

	"github.com/gorilla/websocket"
)

// ListenToAudio joins the audio channel of the meeting and starts listening to the audio stream.
func (c *Client) ListenToAudio() error {

	// Get the STUN and TURN servers
	stunServers, turnServers, err := c.GetStunTurnServers()
	if err != nil {
		return err
	}
	if len(turnServers) == 0 {
		return errors.New("bbb api: No turn servers provided")
	}


	// Make api request to get all information of this meeting (VoiceBridge, CaleeName, UserID, UserName)
	meetings, err := c.API.GetMeetings()
	 if err != nil {
		return err
	}
	meeting := meetings[c.ExternalMeetingID]
	voiceBridge := meeting.VoiceBridge
	caleeName := "GLOBAL_AUDIO_" + strconv.FormatInt(int64(voiceBridge), 10)



	// Connect to the signalling server
	wscon, err := connectToWebSocketSignallingServer(c.WebRTCWSURL, c.SessionToken, c.SessionCookie)
	if err != nil {
		return err
	}

	// Send join message
	err = sendJoinMessage(wscon, c.InternalMeetingID, voiceBridge, caleeName, c.InternalUserID, c.UserName)
	if err != nil {
		return err
	}

	// Read join response
	joinResponse, err := readJoinSDPResponse(wscon)
	if err != nil {
		return err
	}
	sdpOffer := joinResponse.SdpAnswer



	// Create a PeerConnection and set the remote description (sdpAnswer)
	peerConnection, err := createPeerConnection(stunServers, turnServers, sdpOffer)
	if err != nil {
		return err
	}

	// Generate SDP-Offer
	sdpAnswer, err := generateSDPAnswer(peerConnection)
	if err != nil {
		return err
	}



	// Send SDP offer
	err = sendSubscriberSDPAnswer(wscon, voiceBridge, sdpAnswer)
	if err != nil {
		return err
	}


	// Start ping loop
	pingStopChan := pingloop(wscon)
	defer func() {
		pingStopChan <- true
	}()


	// Read SDP answer
	status, err := readStatusResponse(wscon)
	if err != nil {
		return err
	}
	if status.Success != "MEDIA_FLOWING" {
		return errors.New("status response was not successful. Unable to establish webrtc audio connection. ID: " + status.ID + " ,Type: " + status.Type + " ,Success: " + status.Success)
	}

	return nil
}






// STUN and TURN servers
type stunTurns struct {
	StunServers []stunServers `json:"stunServers"`
	TurnServers []turnServers  `json:"turnServers"`
	// RemoteIceCandidates []any `json:"remoteIceCandidates"`
}

type stunServers struct {
	URL string `json:"url"`
}

type turnServers struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	TTL      int    `json:"ttl"`
}

// GetStunTurnServers returns the STUN and TURN servers of the bbb server.
func (c *Client) GetStunTurnServers() ([]stunServers, []turnServers, error) {

	var stunTurns stunTurns

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
		return stunTurns.StunServers, stunTurns.TurnServers, errors.New("bbb api: Couldnt make request. Error: " + err.Error())
	}
	if resp.StatusCode != 200 {
		return stunTurns.StunServers, stunTurns.TurnServers, errors.New("bbb api: Couldnt get stun server address. Server returned: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return stunTurns.StunServers, stunTurns.TurnServers, err
	}

	// Unmarshal json body to stuns struct
	err = json.Unmarshal(body, &stunTurns)
	if err != nil {
		return stunTurns.StunServers, stunTurns.TurnServers, errors.New("bbb api: Couldnt unmarshal server response. Error: " + err.Error())
	}

	return stunTurns.StunServers, stunTurns.TurnServers, nil
}






// Create a PeerConnection
func createPeerConnection(stunServers []stunServers, turnServers []turnServers, sdpOffer string) (*webrtc.PeerConnection, error) {
	// Create ice servers
    iceServers := []webrtc.ICEServer{}

    for _, stun := range stunServers {
        iceServers = append(iceServers, webrtc.ICEServer{
            URLs: []string{stun.URL},
        })
    }

    for _, turn := range turnServers {
        iceServers = append(iceServers, webrtc.ICEServer{
            URLs:     []string{turn.URL},
            Username: turn.Username,
            Credential: turn.Password,
        })
    }

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: iceServers,
	})
	if err != nil {
		return nil, errors.New("failed to create new peer connection: " + err.Error())
	}


	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateFailed {
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
		}
	})


	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  sdpOffer, // This is the SDP offer from the join response
	})
	if err != nil {
		return nil, errors.New("failed to set remote description: " + err.Error())
	}


	// Create an audio transmition
	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio)
	if err != nil {
		return nil, errors.New("failed to add audio transceiver: " + err.Error())
	}

	return peerConnection, nil
}






// Generate a SDP offer and
func generateSDPAnswer(peerConnection *webrtc.PeerConnection) (string, error) {
	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	
	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return "", errors.New("failed to create SDP offer: " + err.Error())
	}


	// Set the local description
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		return "", errors.New("failed to set local description: " + err.Error())
	}


	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// modify the SDP offer for freeswitch
	answer.SDP = rewriteSDP(answer.SDP)

	return answer.SDP, nil
}

// Apply the following transformations for FreeSWITCH
// * Add fake srflx candidate to each media section
// * Add msid to each media section
// * Make bundle first attribute at session level.
func rewriteSDP(in string) string {
	parsed := &sdp.SessionDescription{}
	if err := parsed.Unmarshal([]byte(in)); err != nil {
		panic(err)
	}

	// Reverse global attributes
	for i, j := 0, len(parsed.Attributes)-1; i < j; i, j = i+1, j-1 {
		parsed.Attributes[i], parsed.Attributes[j] = parsed.Attributes[j], parsed.Attributes[i]
	}

	parsed.MediaDescriptions[0].Attributes = append(parsed.MediaDescriptions[0].Attributes, sdp.Attribute{
		Key:   "candidate",
		Value: "79019993 1 udp 1686052607 1.1.1.1 9 typ srflx",
	})

	out, err := parsed.Marshal()
	if err != nil {
		panic(err)
	}

	return string(out)
}






// 
func connectToWebSocketSignallingServer(webrtcwsurl, token string, cookies []*http.Cookie) (*websocket.Conn, error) {
	// Parse url
	u, err := url.Parse(webrtcwsurl + "?sessionToken=" + token)
	if err != nil {
		return nil, errors.New("failed to parse WebRTC WebSocket URL: " + err.Error())
	}
	wsurl := u.String()


	// Create header
	header := http.Header{}

	// Create CookieJar
	coockieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("failed to create cookie jar: " + err.Error())
	}
	tempurl := url.URL(*u)
	tempurl.Scheme = "https"
	coockieJar.SetCookies(&tempurl, cookies)


	// Create dialer
	wsdialer := websocket.DefaultDialer
	wsdialer.Jar = coockieJar


	// Connect to the WebSocket signalling server
	conn, _, err := wsdialer.Dial(wsurl, header)
	if err != nil {
		return nil, errors.New("failed to connect to WebSocket: " + err.Error())
	}
	fmt.Println("Connected to WebSocket server!")
	
	return conn, nil
}


type joinMessage struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Role              string `json:"role"`
	InternalMeetingID string `json:"internalMeetingId"`
	VoiceBridge       int `json:"voiceBridge"`
	CaleeName         string `json:"caleeName"`
	UserID            string `json:"userId"`
	UserName          string `json:"userName"`
}
// Send join message
func sendJoinMessage(wscon *websocket.Conn, internalMeetingID string, voiceBridge int, caleeName, userID, userName string) error {
	// Create join message
	joinMsg := joinMessage{
		ID:                "start",
		Type:              "audio",
		Role:              "recv",
		InternalMeetingID: internalMeetingID,
		VoiceBridge:       voiceBridge,
		CaleeName:         caleeName,
		UserID:            userID,
		UserName:          userName,
	}

	// Marshal join message and send it
	err := wscon.WriteJSON(joinMsg)
	if err != nil {
		return errors.New("failed to send joinMessage: " + err.Error())
	}

	return nil
}

type joinSDPResponse struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Response  string `json:"response"`
	SdpAnswer string `json:"sdpAnswer"`
}
// Read join response
func readJoinSDPResponse(wscon *websocket.Conn) (joinSDPResponse, error) {
	// Read join response
	var joinResp joinSDPResponse
	err := wscon.ReadJSON(&joinResp)
	if err != nil {
		return joinResp, errors.New("failed to read join sdp response: " + err.Error())
	}

	return joinResp, nil
}




type subscriberSDPAnswer struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Role        string `json:"role"`
	VoiceBridge int `json:"voiceBridge"`
	SdpOffer    string `json:"sdpOffer"`
}
// Send SDP offer
func sendSubscriberSDPAnswer(wscon *websocket.Conn, voiceBridge int, answer string) error {
	// Create SDP offer message
	sdpOfferMsg := subscriberSDPAnswer{
		ID:          "subscriberAnswer",
		Type:        "audio",
		Role:        "recv",
		VoiceBridge: voiceBridge,
		SdpOffer:    answer,
	}

	// Marshal SDP offer message and send it
	err := wscon.WriteJSON(sdpOfferMsg)
	if err != nil {
		return errors.New("failed to send SDP offer: " + err.Error())
	}

	return nil
}


type statusResponse struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Success string `json:"success"`
}
// Read status answer
func readStatusResponse(wscon *websocket.Conn) (statusResponse, error) {
	// Read SDP answer
	var sdpRespMsg statusResponse
	err := wscon.ReadJSON(&sdpRespMsg)
	if err != nil {
		return sdpRespMsg, errors.New("failed to read status response: " + err.Error())
	}

	return sdpRespMsg, nil
}





type pingpong struct {
	ID string `json:"id"`
}

// Ping loop
func pingloop(wscon *websocket.Conn) chan bool {
	// Create ping message
	pingMsg := pingpong{
		ID: "ping",
	}

	stopChan := make(chan bool)

	// Start ping loop
	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("Stopping ping loop")
				return
			case <-time.After(15 * time.Second):
				err := wscon.WriteJSON(pingMsg)
				if err != nil {
					fmt.Errorf("failed to send ping message: %v", err)
					return  // Exit the goroutine if there's an error
				}
				fmt.Println("sent ping message!")

				// Set a deadline for reading the next pong message
				wscon.SetReadDeadline(time.Now().Add(5 * time.Second))
				var pongResp pingpong
				err = wscon.ReadJSON(&pongResp)
				if err != nil {
					fmt.Errorf("failed to read pong response: %v", err)
					return  // Exit the goroutine if there's an error
				}

				if pongResp.ID == "pong" {
					fmt.Println("received pong response!")
				} else {
					fmt.Println("received msg, but not pong response!")
				}
			}
		}
	}()

	return stopChan
}






// Connect to the TURN server as a client
func connectToTurnServer(turnServer turnServers) error {
    // Split the URL into host and port
    u, err := url.Parse(turnServer.URL)
    if err != nil {
        return fmt.Errorf("failed to parse TURN server URL: %v", err)
    }

	host, _, err := net.SplitHostPort(u.Opaque)
	if err != nil {
		return fmt.Errorf("failed to parse TURN server URL: %v", err)
	}
    turnServerAddr := u.Opaque


	// Create a local listening socket
    conn, err := net.ListenPacket("udp4", "0.0.0.0:0")
    if err != nil {
        return fmt.Errorf("failed to create local listening socket: %v", err)
    }
    defer conn.Close()


	// Create a Client
    loggerFactory := logging.NewDefaultLoggerFactory()
    loggerFactory.DefaultLogLevel = logging.LogLevelTrace // Setzen Sie das LogLevel auf Trace für detaillierte Protokolle
    cfg := &turn.ClientConfig{
        STUNServerAddr: turnServerAddr,
        TURNServerAddr: turnServerAddr,
        Conn:           conn,
        Username:       turnServer.Username,
        Password:       turnServer.Password,
        Realm:          host, // Nehmen Sie den Host ohne Port als Realm
        LoggerFactory:  loggerFactory,
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

	con, err := client.SendBindingRequest()
	if err != nil {
		return fmt.Errorf("TURN client failed to send binding request: %v", err)
	}

	fmt.Println(con.String())

	
    // // Allocate a relay socket on the TURN server. On success, it
    // // will return a net.PacketConn which represents the remote
    // // socket.
    // relayConn, err := client.Allocate()
    // if err != nil {
    //     return fmt.Errorf("TURN client failed to allocate: %v", err)
    // }
    // defer relayConn.Close()

    // fmt.Printf("Connected to TURN server. Relayed address: %s", relayConn.LocalAddr().String())

    return nil
}