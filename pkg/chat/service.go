package chat

type Service interface {
	Start()
	HasUser(username string) bool
	GetUser(username string) *User
}

type Bot struct {
	Users map[string]*User
}

type Message struct {
	Content string
	From    string
}

type User struct {
	Name           string
	UnReadMessages []*Message
	ReadMessages   []*Message
}

func NewService() Service {
	users := make(map[string]*User)
	return &Bot{
		Users: users,
	}
}

func (b *Bot) Start() {

}

func (b *Bot) HasUser(username string) bool {
	return true
}

func (b *Bot) GetUser(username string) *User {
	return nil
}
