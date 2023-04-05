package fcm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"main.go/databasehandler"
)

func SendNotification( user_id string, firstName string, middleName string) {
	sendbody := fmt.Sprintf("%v %v is out of the location radius set", firstName, middleName)
	url := "https://fcm.googleapis.com/fcm/send"
	fmt.Println("URL:>", url)

	jsonStr := fmt.Sprintf(`{
		"to": "/topics/%v",
		"notification": {
		  "body": "%v",
		  "OrganizationId": "2",
		  "content_available": true,
		  "priority": "high",
		  "subtitle": "Elementary School",
		  "title": "Geofence Breached"
		},
		"data": {
		  "priority": "high",
		  "user": "charles",
		  "content_available": true,
		  "bodyText": "New Announcement assigned",
		  "organization": "Elementary school"
		}
	}`,user_id, sendbody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", databasehandler.GoDotEnvVariable("FCMSERVERKEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
