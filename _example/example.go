package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"api"

	bot "github.com/ITLab-CC/bigbluebutton-bot"
)

type config struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
}

func readConfig(file string) config {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	// we initialize config
	var conf config
	// we unmarshal our byteArray which contains our jsonFile's content into conf
	json.Unmarshal([]byte(byteValue), &conf)

	return conf
}

func main() {

	conf := readConfig("config.json")

	bbbapi, err := api.NewRequest(conf.Url, conf.Secret, api.SHA256)
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

	time.Sleep(5 * time.Second)



	fmt.Println("-----------------------------------------------")



	client, err := bot.NewClient("bbb6.ccita.de", conf.Secret)
	if err != nil {
		panic(err)
	}

	client.OnStatus(func(status bot.Status) {
		fmt.Printf("Bot status: %s\n", status)
	})

	fmt.Println("Bot joins " + newmeeting.MeetingName + " as moderator:")
	err = client.Join(newmeeting.MeetingID, "Bot", true)
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)



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