package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":8096", nil)

}

type GcpData struct {
	Incident struct {
		IncidentID    string `json:"incident_id"`
		ResourceID    string `json:"resource_id"`
		ResourceName  string `json:"resource_name"`
		State         string `json:"state"`
		StartedAt     *int64 `json:"started_at"`
		EndedAt       *int64 `json:"ended_at"`
		PolicyName    string `json:"policy_name"`
		ConditionName string `json:"condition_name"`
		URL           string `json:"url"`
		Summary       string `json:"summary"`
	} `json:"incident"`
}

type Notify struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	var (
		input  *GcpData
		output Notify
		err    error
	)

	url := "https://api.telegram.org/bot725691005:AAFGVkxe-mdDrLrlpeyPNG5c3dqjtD0dhPg/sendMessage"
	output.ChatID = -300312399

	session := r.Header.Get("Authorization")
	if session != "Basic aW5zbG90OnF3ZTEyMw==" {
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	if err = json.Unmarshal(buf, &input); err != nil {
		log.Printf("😈  Json Decode Error ---> %+v \n", err)
		output.Text = string(buf)
		SendMessage("POST", url, output)
		return
	}

	if input == nil {
		log.Printf("😈  輸入資料為null \n")
		return
	}
	var (
		startAt int64
		endAt   int64
	)

	if input.Incident.StartedAt != nil {
		startAt = *input.Incident.StartedAt
	}
	if input.Incident.EndedAt != nil {
		endAt = *input.Incident.StartedAt
	}

	output.Text = fmt.Sprintf(`
	⭐  %s  ：  %s  ⭐

	⏰  時間： %v  ～  %v

	✍  網址： %v

	`,
		input.Incident.ConditionName, input.Incident.Summary,
		time.Unix(startAt, 0).UTC().Add(time.Hour*8).Format("2006-01-02 15:04:05"), time.Unix(endAt, 0).UTC().Add(time.Hour*8).Format("2006-01-02 15:04:05"),
		input.Incident.URL,
	)

	SendMessage("POST", url, output)

}

// SendMessage 發送訊息
func SendMessage(method, url string, output Notify) {

	outputData, err := json.Marshal(output)
	if err != nil {
		log.Printf("😈  Json Marshal Error ---> %v \n", err)
		return
	}

	payload := strings.NewReader(string(outputData))

	var (
		req *http.Request
	)
	req, err = http.NewRequest("POST", url, payload)
	if err != nil {
		log.Printf("😈  NewRequest Error ---> %v \n", err)
		return
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Basic aW5zbG90OnF3ZTEyMw==")
	req.Header.Add("cache-control", "no-cache")

	_, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		log.Printf("😈  Do Error ---> %v \n", err)
		return
	}

	return
}
