package gogitter

import (
  "path"
  "bytes"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "strconv"
)

const streamingDomain = "https://stream.gitter.im/"
const apiDomain = "https://api.gitter.im/"
const version = "v1"

func GetMessageStream(token string, roomId string) (chan GitterMessage, func()) {
  url := streamingDomain + path.Join(version,
    "rooms", roomId, "chatMessages")
  return GitterEventStream(url, token)
}

func GetSendMessageStream(token string, roomId string) chan string {
  url := apiDomain + path.Join(version,
    "rooms", roomId, "chatMessages")
  in := make(chan string)
  go func() {
    for text := range in {
      client := &http.Client{}
      payload, _ := json.Marshal(GitterSendMessage{text})
      req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
      req.Header.Add("Authorization", "Bearer " + token)
      req.Header.Add("Content-Type", "application/json")
      req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
      resp, _ := client.Do(req)
      data, _ := ioutil.ReadAll(resp.Body)
      message := GitterMessage{}
      json.Unmarshal(data, &message)
    }
  }()
  return in
}
