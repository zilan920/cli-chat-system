package chat

import "errors"

type Service interface {
	Start()
	CallCmd(cmdName string, args []string) (error, string)
	Login(username string) (error, string)
	SendMsg(to, msg string) (error, string)
	CheckMsg() (error, []string)
}

type Bot struct {
	Users       map[string]*User
	CurrentUser *User
}

type User struct {
	Name     string
	Messages []string
}

type Message struct {
	Content string
	From    string
}

func NewService() Service {
	users := make(map[string]*User)
	return &Bot{
		Users:       users,
		CurrentUser: nil,
	}
}

func (b *Bot) Start() {

}

func (b *Bot) Login(username string) (error, string) {
	if !b.hasUser(username) {
		messages := make([]string, 10)
		user := &User{
			Name:     username,
			Messages: messages,
		}
		b.Users[username] = user
		b.CurrentUser = user
	}
	return nil, ""
}

func (b *Bot) SendMsg(to, msg string) (error, string) {
	if b.CurrentUser == nil {
		return errors.New("please login First"), ""
	}
	user, ok := b.Users[to]
	if !ok {
		return errors.New("target user not found"), ""
	}
	user.Messages = append(user.Messages, msg)
	return nil, "msg sent"
}

func (b *Bot) CheckMsg() (error, []string) {
	if b.CurrentUser == nil {
		return errors.New("please login First"), []string{}
	}
	messages := b.CurrentUser.Messages
	return nil, messages
}

func (b *Bot) hasUser(username string) bool {
	_, ok := b.Users[username]
	return ok
}

func (b *Bot) CallCmd(cmdName string, args []string) (error, string) {
	switch cmdName {
	case "login":
		err, output := b.Login(args[0])
		return err, output
	case "send":
		err, output := b.SendMsg(args[0], args[1])
		return err, output
	case "read":
		err, outputs := b.CheckMsg()
		output := outputs[0]
		return err, output
	case "reply":
		err, output := b.SendMsg("", args[0])
		return err, output
	case "forward":
		err, output := b.SendMsg("", args[0])
		return err, output
	case "broadcast":
		err, output := b.SendMsg("", args[0])
		return err, output
	default:
		return errors.New("invalid cmd"), ""
	}
}
