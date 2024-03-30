package onebot

import "encoding/json"

type Event struct {
	PostType string `json:"post_type"`
}

func UnmarshalEvent(data []byte) (*Event, error) {
	var event *Event = &Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}

type Ret struct {
	RetCode int    `json:"retcode"`
	Status  string `json:"status"`
	Echo    string `json:"echo"`
}

type User struct {
	UserId   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

type GetFriendListRet struct {
	Ret
	Data []User `json:"data"`
}
