package notifier

type Zoom struct {
	URL   string
	Token string
}

type ZoomPayload struct {
	RobotJid  string `json:"robot_jid"`
	ToJid     string `json:"to_jid"`
	AccountId string `json:"account_id"`
	Content  ZoomContent `json:"content"`
}

type ZoomContent struct {

}
