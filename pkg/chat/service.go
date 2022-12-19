package chat

import "fmt"

type Service interface {
	Start()
	HasUser(username string) bool
	GetUser(username string) *User
	CreateUser(username string) *User
	RegisterTenant(username string, tenant ServiceTenant) error
	UnRegisterTenant(username string) error
	SendMsg(target *User, msg, from string) error
	Broadcast(msg string) error
}

type Bot struct {
	Users   map[string]*User
	Tenants map[string]ServiceTenant
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
		Users:   users,
		Tenants: map[string]ServiceTenant{},
	}
}

func (b *Bot) Start() {
}

func (b *Bot) HasUser(username string) bool {
	_, ok := b.Users[username]
	return ok
}

func (b *Bot) GetUser(username string) *User {
	user, ok := b.Users[username]
	if ok {
		return user
	}
	return nil
}

func (b *Bot) CreateUser(username string) *User {
	if b.HasUser(username) {
		return b.GetUser(username)
	}
	urm := make([]*Message, 0)
	rm := make([]*Message, 0)
	user := User{
		Name:           username,
		UnReadMessages: urm,
		ReadMessages:   rm,
	}
	b.Users[username] = &user
	return &user
}

func (b *Bot) RegisterTenant(username string, tenant ServiceTenant) error {
	b.Tenants[username] = tenant
	return nil
}

func (b *Bot) UnRegisterTenant(username string) error {
	delete(b.Tenants, username)
	return nil
}

func (b *Bot) SendMsg(target *User, msg, from string) error {
	target.UnReadMessages = append(target.UnReadMessages, &Message{
		Content: msg,
		From:    from,
	})
	if tenant, ok := b.Tenants[target.Name]; ok {
		tenant.Push(fmt.Sprintf("you have a new message from %s \n", from))
	}
	return nil
}

func (b *Bot) Broadcast(msg string) error {
	return nil
}
