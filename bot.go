package bot

import (
	api "api"
	"errors"

	ddp "ddp"
)

type Status string
const (
	DISCONNECTED 	Status = "disconected"
	CONNECTING   	Status = "connecting"
	CONNECTED 		Status = "connected"
)


// Client represents a BigBlueButton client connection. The BigBlueButton client establish a BigBlueButton
// session and acts as a message pump for other tools.
type Client struct {
	// connectionStatus is the current connection status of the client
	connectionStatus 	Status

	// BBB-urls the client is connected to
	clientURL			string
	clientWSURL			string
	apiURL				string
	apiSecret			string
	// to make api requests to the BBB-server
	API 				*api.ApiRequest

	ddpClient 			*ddp.Client

	event 				*event
}

func NewClient(clientURL string, clientWSURL string, apiURL string, apiSecret string) (*Client, error) {
	api, err := api.NewRequest(apiURL, apiSecret, api.SHA256)
	if (err != nil) {
		return nil, err
	}
	
	ddpClient := ddp.NewClient(clientWSURL, clientURL)

	c := &Client{
		connectionStatus: 	DISCONNECTED,

		clientURL:			clientURL,
		clientWSURL:		clientWSURL,
		apiURL:				apiURL,
		apiSecret:			apiSecret,

		ddpClient:			ddpClient,

		API: 				api,

		event: 				nil,
	}

	e := &event{
		client: c,
	}

	c.event = e
	c.ddpClient.AddStatusListener(e)

	return c, nil
}

func (c *Client) Join(meetingID string, userName string, moderator bool) error {
	_, _, internalUserID, authToken, _, internalMeetingID, err := c.API.Join(meetingID, userName, moderator)
	if err != nil {
		return err
	}

	err = c.ddpClient.Connect()
	if err != nil {
		return err
	}

	err = c.ddpClient.Sub("meteor_autoupdate_clientVersions")
	if err != nil {
		return err
	}

	settings := `[
		{
			"application": {
			"animations": false,
			"chatAudioAlerts": false,
			"chatPushAlerts": false,
			"userJoinAudioAlerts": false,
			"userJoinPushAlerts": false,
			"userLeaveAudioAlerts": false,
			"userLeavePushAlerts": false,
			"raiseHandAudioAlerts": false,
			"raiseHandPushAlerts": false,
			"guestWaitingAudioAlerts": false,
			"guestWaitingPushAlerts": false,
			"paginationEnabled": false,
			"pushLayoutToEveryone": false,
			"fallbackLocale": "en",
			"overrideLocale": null,
			"locale": "en"
			},
			"audio": {
			"inputDeviceId": "undefined",
			"outputDeviceId": "undefined"
			},
			"dataSaving": {
			"viewParticipantsWebcams": true,
			"viewScreenshare": true
			}
		}
		]
	}`

	_, err = c.ddpClient.Call("userChangedLocalSettings", settings)
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