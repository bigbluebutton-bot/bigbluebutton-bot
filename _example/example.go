package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	api "github.com/ITLab-CC/bigbluebutton-bot/api"

	bot "github.com/ITLab-CC/bigbluebutton-bot"

	bbb "github.com/ITLab-CC/bigbluebutton-bot/bbb"
)

type configAPI struct {
	URL    string  `json:"url"`
	Secret string  `json:"secret"`
	SHA    api.SHA `json:"sha"`
}

type configClient struct {
	URL string `json:"url"`
	WS  string `json:"ws"`
}

type configPad struct {
	URL string `json:"url"`
	WS  string `json:"ws"`
}

type configWebRTC struct {
	WS string `json:"ws"`
}

type configBBB struct {
	API    configAPI	`json:"api"`
	Client configClient	`json:"client"`
	Pad configPad		`json:"pad"`
	WebRTC configWebRTC `json:"webrtc"`
}

type config struct {
	BBB configBBB `json:"bbb"`
}

func readConfig(file string) config {
	// Try to read from env
	conf := config {
		BBB: configBBB{
			API: configAPI{
				URL: os.Getenv("BBB_API_URL"),
				Secret: os.Getenv("BBB_API_SECRET"),
				SHA: api.SHA(os.Getenv("BBB_API_SHA")),
			},
			Client: configClient{
				URL: os.Getenv("BBB_CLIENT_URL"),
				WS: os.Getenv("BBB_CLIENT_WS"),
			},
			Pad: configPad{
				URL: os.Getenv("BBB_PAD_URL"),
				WS: os.Getenv("BBB_PAD_WS"),
			},
			WebRTC: configWebRTC{
				WS: os.Getenv("BBB_WEBRTC_WS"),
			},
		},
	}

	if (conf.BBB.API.URL != "" && conf.BBB.API.Secret != "" && conf.BBB.API.SHA != "" &&
		conf.BBB.Client.URL != "" && conf.BBB.Client.WS != "" &&
		conf.BBB.Pad.URL != "" && conf.BBB.Pad.WS != ""){
		fmt.Println("Using env variables for config")
		return conf
	}

	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if (err != nil) {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if(err != nil) {
		panic(err)
	}
	// we unmarshal our byteArray which contains our jsonFile's content into conf
	json.Unmarshal([]byte(byteValue), &conf) 

	return conf
}

func main() {

	conf := readConfig("config.json")

	bbbapi, err := api.NewRequest(conf.BBB.API.URL, conf.BBB.API.Secret, conf.BBB.API.SHA)
	if err != nil {
		panic(err)
	}


	//API-Requests
	newmeeting, err := bbbapi.CreateMeeting("name", "meetingID", "attendeePW", "moderatorPW", "welcome text", false, false, false, 12345)
	if err != nil {
		panic(err)
	}
	fmt.Printf("New meeting \"%s\" was created.\n", newmeeting.MeetingName)



	fmt.Println("-----------------------------------------------")



	fmt.Println("All meetings:")
	meetings, err := bbbapi.GetMeetings()
	if err != nil {
		panic(err)
	}
	for _, meeting := range meetings {
		fmt.Print(meeting.MeetingName + ": ")
		fmt.Println(bbbapi.IsMeetingRunning(meeting.MeetingID))
	}



	fmt.Println("-----------------------------------------------")



	url, err := bbbapi.JoinGetURL(newmeeting.MeetingID, "TestUser", true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Moderator join url: " + url)

	time.Sleep(1 * time.Second)



	fmt.Println("-----------------------------------------------")



	client, err := bot.NewClient(conf.BBB.Client.URL, conf.BBB.Client.WS, conf.BBB.Pad.URL, conf.BBB.Pad.WS, conf.BBB.API.URL, conf.BBB.API.Secret, conf.BBB.WebRTC.WS)
	if err != nil {
		panic(err)
	}

	client.OnStatus(func(status bot.StatusType) {
		fmt.Printf("Bot status: %s\n", status)
	})

	fmt.Println("Bot joins " + newmeeting.MeetingName + " as moderator:")
	err = client.Join(newmeeting.MeetingID, "Bot", true)
	if err != nil {
		panic(err)
	}

	err = client.OnGroupChatMsg(func(msg bbb.Message) {

		fmt.Println("[" + msg.SenderName + "]: " + msg.Message)

		if(msg.Sender != client.InternalUserID) {
			if(msg.Message == "ping") {
				fmt.Println("Sending pong")
				client.SendChatMsg("pong", msg.ChatId)
			}
		}
	})
	if err != nil {
		panic(err)
	}

	// enCapture, err := client.CreateCapture("en")
	// if err != nil {
	// 	panic(err)
	// }

	// time.Sleep(2 * time.Second)

	// err = enCapture.SendText("Hello")
	// if err != nil {
	// 	panic(err)
	// }

	// time.Sleep(2 * time.Second)

	// err = enCapture.SendText(" world")
	// if err != nil {
	// 	panic(err)
	// }

	// time.Sleep(2000 * time.Second)


	err = client.ListenToAudio()
	if err != nil {
		panic(err)
	}

	

	fmt.Println("Bot leaves " + newmeeting.MeetingName)
	err = client.Leave()
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	

	fmt.Println("-----------------------------------------------")



	endedmeeting, err := bbbapi.EndMeeting(newmeeting.MeetingID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Meeting \"%s\" was ended.\n", endedmeeting.MeetingName)



	fmt.Println("-----------------------------------------------")



	fmt.Println("All meetings:")
	meetings, err = bbbapi.GetMeetings()
	if err != nil {
		panic(err)
	}
	for _, meeting := range meetings {
		fmt.Println(meeting.MeetingName)
	}
}
