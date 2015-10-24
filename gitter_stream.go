package gogitter

import (
  "net/http"
  "strings"
  "encoding/json"
)

func GitterEventStream(url, token string) (chan GitterMessage, func()) {
  out := make(chan GitterMessage)
  stop := make(chan bool)
  go func() {
    client := &http.Client{}
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", "Bearer " + token)
    resp, _ := client.Do(req)
    gitterParseEventStream(resp, out, stop)
    resp.Body.Close()
    close(out)
  }()

  return out, func () {
    stop <- true
  }
}

func gitterStreamReciever(resp *http.Response, out chan string, stop chan bool) {
  var buf []byte
  flag := false
  tube := make(chan string)
  for {
    buf = make([]byte, 1024)
    go func() {
      resp.Body.Read(buf)
      tube <- string(buf)
    }()
    select {
    case msg := <-tube:
      out <- msg
    case <-stop:
      flag = true
      close(out)
    }
    if flag {
      break
    }
  }
}

func gitterParseEventStream(resp *http.Response, out chan GitterMessage, stop chan bool) {
  rin := make(chan string)
  go gitterStreamReciever(resp, rin, stop)
  for resStr := range rin {
    parts := strings.Split(resStr, "\n")
    for i, part := range parts {
      if i < len(parts) - 1 && len(part) > 1 {
        message := GitterMessage{}
        backtob := []byte(strings.TrimSpace(part))
        json.Unmarshal(backtob, &message)
        out <- message
      }
    }
  }
}
