package chat

import (
	"errors"
	"fmt"
)

type ServiceTenant interface {
	CallCmd(cmdName string, args []string) (error, string)
	Login(username string) (error, string)
	SendMsg(to, msg string) (error, string)
	CheckMsg() (error, *Message)
	CheckNewMsg() (error, *Message)
}

type TenantContext struct {
	Message *Message
}

type Tenant struct {
	svc         Service
	CurrentUser *User
	Context     *TenantContext
}

func NewTenant(svc Service) ServiceTenant {
	return &Tenant{
		svc: svc,
	}
}

func (t *Tenant) Login(username string) (error, string) {
	if !t.svc.HasUser(username) {
		unreadMessages := make([]*Message, 0)
		readMessages := make([]*Message, 0)
		user := &User{
			Name:           username,
			UnReadMessages: unreadMessages,
			ReadMessages:   readMessages,
		}
		t.Users[username] = user
		t.CurrentUser = user
	} else {
		user, ok := t.Users[username]
		if !ok {
			return errors.New("target user not found"), ""
		}
		t.CurrentUser = user
	}
	t.Context = nil
	return nil, fmt.Sprintf("%s logged in, %d new messages", username, len(b.CurrentUser.UnReadMessages))
}

func (t *Tenant) SendMsg(to, msg string) (error, string) {
	if t.CurrentUser == nil {
		return errors.New("please login First"), ""
	}
	user, ok := t.Users[to]
	if !ok {
		return errors.New("target user not found"), ""
	}
	user.UnReadMessages = append(user.UnReadMessages, &Message{
		Content: msg,
		From:    b.CurrentUser.Name,
	})
	return nil, "msg sent"
}

func (t *Tenant) CheckMsg() (error, *Message) {
	if t.CurrentUser == nil {
		return errors.New("please login First"), &Message{}
	}
	messages := t.CurrentUser.UnReadMessages
	if len(messages) > 0 {
		message := messages[len(messages)-1]
		return nil, message
	}
	return nil, nil
}

func (t *Tenant) CheckNewMsg() (error, *Message) {
	if t.CurrentUser == nil {
		return errors.New("please login First"), &Message{}
	}
	var message *Message
	messages := t.CurrentUser.UnReadMessages
	unReadCount := len(messages)
	if unReadCount > 0 {
		message = messages[unReadCount-1]
		t.CurrentUser.ReadMessages = append(t.CurrentUser.ReadMessages, message)
		t.CurrentUser.UnReadMessages = t.CurrentUser.UnReadMessages[0 : unReadCount-1]
		t.Context = &BotContext{
			Message: message,
		}
	}
	return nil, message
}

func (t *Tenant) CallCmd(cmdName string, args []string) (error, string) {
	switch cmdName {
	case "login":
		err, output := t.Login(args[0])
		return err, output
	case "send":
		err, output := t.SendMsg(args[0], args[1])
		return err, output
	case "read":
		err, message := t.CheckNewMsg()
		if message == nil {
			return err, fmt.Sprintf("you don't have any new messages")
		}
		return err, fmt.Sprintf("from %s: %s", message.From, message.Content)
	case "reply":
		if t.Context == nil {
			return errors.New("no reply target"), ""
		}
		err, output := t.SendMsg(t.Context.Message.From, args[0])
		return err, output
	case "forward":
		if t.Context == nil {
			return errors.New("no forward target"), ""
		}
		err, output := t.SendMsg(args[0], t.Context.Message.Content)
		return err, output
	case "broadcast":
		err, output := t.SendMsg("", args[0])
		return err, output
	default:
		return errors.New("invalid cmd"), ""
	}
}
