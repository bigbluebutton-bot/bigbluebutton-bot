package bot

import (
	api "github.com/ITLab-CC/bigbluebutton-bot/api"
	"errors"
	"net/http"
	"time"

	ddp "github.com/gopackage/ddp"

	bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"
)

type StatusType string

const (
	DISCONNECTING StatusType = "disconnecting"
	DISCONNECTED  StatusType = "disconnected"
	CONNECTING    StatusType = "connecting"
	CONNECTED     StatusType = "connected"
	RECONNECTING  StatusType = "reconnecting"
)

// Client represents a BigBlueButton client connection. The BigBlueButton client establish a BigBlueButton
// session and acts as a message pump for other tools.
type Client struct {
	// Status is the current connection status of the client
	Status StatusType

	// BBB-urls the client is connected to
	ClientURL   string
	ClientWSURL string
	PadURL      string
	PadWSURL    string
	ApiURL      string
	apiSecret   string
	// to make api requests to the BBB-server
	API *api.ApiRequest

	ddpClient *ddp.Client

	// events will store all the functions executed on certain events. (events["OnStatus"][]func(StatusType))
	events          map[string][]interface{}
	ddpEventHandler *ddpEventHandler

	// after join there are the following informations
	JoinURL           string
	SessionCookie     []*http.Cookie
	InternalUserID    string
	AuthToken         string
	SessionToken      string
	InternalMeetingID string
}

func NewClient(clientURL string, clientWSURL string, padURL string, padWSURL string, apiURL string, apiSecret string) (*Client, error) {
	api, err := api.NewRequest(apiURL, apiSecret, api.SHA256)
	if err != nil {
		return nil, err
	}

	ddpClient := ddp.NewClient(clientWSURL, clientURL)

	c := &Client{
		Status: DISCONNECTED,

		ClientURL:   clientURL,
		ClientWSURL: clientWSURL,
		PadURL:      padURL,
		PadWSURL:    padWSURL,
		ApiURL:      apiURL,
		apiSecret:   apiSecret,

		ddpClient: ddpClient,

		API: api,

		events:          make(map[string][]interface{}),
		ddpEventHandler: nil,
	}

	c.ddpEventHandler = &ddpEventHandler{
		client:  c,
		updater: make(map[string][]updaterfunc),
	}

	return c, nil
}

// Join a meeting
func (c *Client) Join(meetingID string, userName string, moderator bool) error {
	joinURL, coockie, internalUserID, authToken, sessionToken, internalMeetingID, err := c.API.Join(meetingID, userName, moderator)
	if err != nil {
		return err
	}
	c.JoinURL = joinURL
	c.SessionCookie = coockie
	c.InternalUserID = internalUserID
	c.AuthToken = authToken
	c.SessionToken = sessionToken
	c.InternalMeetingID = internalMeetingID

	// Connect to the DDP server
	if err = c.ddpConnect(); err != nil {
		return err
	}

	// Subscribe to the current user
	if err = c.ddpSubscribe(bbb.CurrentUser, nil); err != nil {
		return err
	}

	// Call the validateAuthToken method with the userID, authToken, and userName
	_, err = c.ddpCall(bbb.ValidateAuthTokenCall, internalMeetingID, internalUserID, authToken, internalUserID)
	if err != nil {
		return errors.New("could not validateAuthToken")
	}

	return nil
}

// Leave the joined meeting
func (c *Client) Leave() error {
	// If not connected, return an error
	if c.Status != CONNECTED {
		// If is connecting retry 5 times
		if c.Status == CONNECTING {
			i := 0
			for i < 5 {
				if c.Status == CONNECTED {
					c.Leave()
				}
				time.Sleep(time.Second * 1)
				i += 1
			}
		}
		return errors.New("Client is in no meeting. First Join a meeting with: client.Join(meetingID string, userName string, moderator bool)")
	}

	c.ddpCall(bbb.UserLeftMeetingCall)
	c.ddpCall(bbb.SetExitReasonCall, "logout")
	// c.ddpClient.UnSubscribe("from all subs")

	c.ddpDisconnect()

	c.ddpClient = nil

	return nil
}
