package chat

import (
	"errors"
	"fmt"
)

type Service interface {
	Start()
	CallCmd(cmdName string, args []string) (error, string)
	Login(username string) (error, string)
	SendMsg(to, msg string) (error, string)
	CheckMsg() (error, *Message)
	CheckNewMsg() (error, *Message)
}
type BotContext struct {
	Message *Message
}

type Bot struct {
	Users       map[string]*User
	CurrentUser *User
	Context     *BotContext
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
		Users:       users,
		CurrentUser: nil,
	}
}

func (b *Bot) Start() {

}

func (b *Bot) Login(username string) (error, string) {
	if !b.hasUser(username) {
		unreadMessages := make([]*Message, 0)
		readMessages := make([]*Message, 0)
		user := &User{
			Name:           username,
			UnReadMessages: unreadMessages,
			ReadMessages:   readMessages,
		}
		b.Users[username] = user
		b.CurrentUser = user
	} else {
		user, ok := b.Users[username]
		if !ok {
			return errors.New("target user not found"), ""
		}
		b.CurrentUser = user
	}
	b.Context = nil
	return nil, fmt.Sprintf("%s logged in, %d new messages", username, len(b.CurrentUser.UnReadMessages))
}

func (b *Bot) SendMsg(to, msg string) (error, string) {
	if b.CurrentUser == nil {
		return errors.New("please login First"), ""
	}
	user, ok := b.Users[to]
	if !ok {
		return errors.New("target user not found"), ""
	}
	user.UnReadMessages = append(user.UnReadMessages, &Message{
		Content: msg,
		From:    b.CurrentUser.Name,
	})
	return nil, "msg sent"
}

func (b *Bot) CheckMsg() (error, *Message) {
	if b.CurrentUser == nil {
		return errors.New("please login First"), &Message{}
	}
	messages := b.CurrentUser.UnReadMessages
	if len(messages) > 0 {
		message := messages[len(messages)-1]
		return nil, message
	}
	return nil, nil
}

func (b *Bot) CheckNewMsg() (error, *Message) {
	if b.CurrentUser == nil {
		return errors.New("please login First"), &Message{}
	}
	var message *Message
	messages := b.CurrentUser.UnReadMessages
	unReadCount := len(messages)
	if unReadCount > 0 {
		message = messages[unReadCount-1]
		b.CurrentUser.ReadMessages = append(b.CurrentUser.ReadMessages, message)
		b.CurrentUser.UnReadMessages = b.CurrentUser.UnReadMessages[0 : unReadCount-1]
		b.Context = &BotContext{
			Message: message,
		}
	}
	return nil, message
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
		err, message := b.CheckNewMsg()
		if message == nil {
			return err, fmt.Sprintf("you don't have any new messages")
		}
		return err, fmt.Sprintf("from %s: %s", message.From, message.Content)
	case "reply":
		if b.Context == nil {
			return errors.New("no reply target"), ""
		}
		err, output := b.SendMsg(b.Context.Message.From, args[0])
		return err, output
	case "forward":
		if b.Context == nil {
			return errors.New("no forward target"), ""
		}
		err, output := b.SendMsg(args[0], b.Context.Message.Content)
		return err, output
	case "broadcast":
		err, output := b.SendMsg("", args[0])
		return err, output
	default:
		return errors.New("invalid cmd"), ""
	}
}
