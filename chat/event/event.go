package event

type MessageType string

const (
	// ENTER 입장
	MessageTypeEnter MessageType = "ENTER"
	// LEAVE 퇴장
	MessageTypeLeave MessageType = "LEAVE"
	// CHAT 메시지
	MessageTypeChat MessageType = "CHAT"
)

type ChatEvent struct {
	UserID      string      `json:"userId"`
	Nickname    string      `json:"nickname"`
	RoomID      string      `json:"roomId"`
	Message     string      `json:"message"`
	MessageType MessageType `json:"messageType"`
}
