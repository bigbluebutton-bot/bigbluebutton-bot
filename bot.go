package bot

import (
	api "api"
	"errors"
	"time"

	ddp "ddp"
)

type StatusType string
const (
	DISCONNECTING 	StatusType = "disconnecting"
	DISCONNECTED 	StatusType = "disconnected"
	CONNECTING   	StatusType = "connecting"
	CONNECTED 		StatusType = "connected"
	RECONNECTING 	StatusType = "reconnecting"
)

// This is for all events that are in "event_....go" files
type event struct {
	client *Client
}


// Client represents a BigBlueButton client connection. The BigBlueButton client establish a BigBlueButton
// session and acts as a message pump for other tools.
type Client struct {
	// Status is the current connection status of the client
	Status 	StatusType

	// BBB-urls the client is connected to
	ClientURL			string
	ClientWSURL			string
	ApiURL				string
	apiSecret			string
	// to make api requests to the BBB-server
	API 				*api.ApiRequest

	ddpClient 			*ddp.Client

	// events will store all the functions executed on certain events. (events["OnStatus"][]func(StatusType))
	events 				map[string][]interface{}
}

func NewClient(clientURL string, clientWSURL string, apiURL string, apiSecret string) (*Client, error) {
	api, err := api.NewRequest(apiURL, apiSecret, api.SHA256)
	if (err != nil) {
		return nil, err
	}
	
	ddpClient := ddp.NewClient(clientWSURL, clientURL)

	c := &Client{
		Status: 	DISCONNECTED,

		ClientURL:			clientURL,
		ClientWSURL:		clientWSURL,
		ApiURL:				apiURL,
		apiSecret:			apiSecret,

		ddpClient:			ddpClient,

		API: 				api,

		events: 			make(map[string][]interface{}),
	}

	return c, nil
}

// Join a meeting
func (c *Client) Join(meetingID string, userName string, moderator bool) error {
	_, _, internalUserID, authToken, _, internalMeetingID, err := c.API.Join(meetingID, userName, moderator)
	if err != nil {
		return err
	}

	err = c.ddpClient.Connect()
	if err != nil {
		return err
	}

	err = c.ddpClient.Sub("current-user")
	if err != nil {
		return errors.New("could sub current-user")
	}

	// Call the validateAuthToken method with the userID, authToken, and userName
	_, err = c.ddpClient.Call("validateAuthToken", internalMeetingID, internalUserID, authToken, internalUserID)
	if err != nil {
		return errors.New("could not validateAuthToken")
	}

	return nil
}

// Leave the joined meeting
func (c *Client) Leave() error {
	// If not connected, return an error
	if(c.Status != CONNECTED) {
		// If is connecting retry 5 times
		if(c.Status == CONNECTING) {
			i := 0
			for(i < 5) {
				if(c.Status == CONNECTED) {
					c.Leave()
				}
				time.Sleep(time.Second * 1)
				i += 1
			}
		}
		return errors.New("Client is in no meeting. First Join a meeting with: client.Join(meetingID string, userName string, moderator bool)")
	}

	c.ddpClient.Call("userLeftMeeting")
	c.ddpClient.Call("setExitReason", "logout")
	// c.ddpClient.UnSubscribe("from all subs")

	c.ddpClient.Close()

	c.ddpClient = nil

	return nil
}