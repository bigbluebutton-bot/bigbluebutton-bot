package bot

import (
	api "api"
	"errors"

	ddp "github.com/gopackage/ddp"
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
	connectionStatus Status
	// // statusListeners will be informed when the connection status of the client changes
	// statusListeners []statusListener

	// BBB-url the client is connected to
	domain string
	// to make api requests to the BBB-server
	API *api.ApiRequest

	ddpClient *ddp.Client

	eventsOnStatus []statusListener
}

func NewClient(domain string, secret string) (*Client, error) {
	api, err := api.NewRequest("https://" + domain + "/bigbluebutton/api", secret, api.SHA256)
	if (err != nil) {
		return nil, err
	}
	
	c := &Client{
		connectionStatus: 	DISCONNECTED,
		domain: 			domain,
		API: 				api,

		eventsOnStatus: 	nil,
	}

	c.ddpClient = ddp.NewClient("wss://" + domain + "/html5client/websocket", "https://" + domain + "/html5client/")
	c.ddpClient.AddStatusListener(c)

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