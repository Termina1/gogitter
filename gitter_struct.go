package gogitter

import "time"

type GitterMessage struct {
  Id string
  Text string
  Html string
  Sent time.Time
  FromUser GitterUser
  Urls []GitterUrl
  Unread bool
  ReadBy int
  V int
  Mentions []GitterMention
}

type GitterUser struct {
  Id string
  Username string
  DisplayName string
  Url string
  AvatarUrlMedium string
  AvatarUrlSmall string
  Gv int
}

type GitterUrl struct {
  Url string
}

type GitterMention struct {
  ScreenName string
  UserIds []string
  UserId string
  group bool
}

type GitterSendMessage struct {
  Text string `json:"text"`
}
