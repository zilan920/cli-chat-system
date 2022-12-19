package chat

import (
	"errors"
	"fmt"
	"net"
)

type ServiceTenant interface {
	CallCmd(cmdName string, args []string) (error, string)
	Login(username string) (error, string)
	SendMsg(to, msg string) (error, string)
	CheckMsg() (error, *Message)
	CheckNewMsg() (error, *Message)
	Bye()
}

type TenantContext struct {
	Message *Message
}

type Tenant struct {
	fd          int
	svc         Service
	conn        net.Conn
	CurrentUser *User
	Context     *TenantContext
}

func NewTenant(fd int, svc Service, conn net.Conn) ServiceTenant {
	return &Tenant{
		fd:   fd,
		svc:  svc,
		conn: conn,
	}
}

func (t *Tenant) Login(username string) (error, string) {
	if !t.svc.HasUser(username) {
		user := t.svc.CreateUser(username)
		t.CurrentUser = user
	} else {
		user := t.svc.GetUser(username)
		t.CurrentUser = user
	}
	t.Context = nil
	_ = t.svc.RegisterTenant(username, t)
	return nil, fmt.Sprintf("%s logged in, %d new messages", username, len(t.CurrentUser.UnReadMessages))
}

func (t *Tenant) SendMsg(to, msg string) (error, string) {
	if t.CurrentUser == nil {
		return errors.New("please login First"), ""
	}
	if !t.svc.HasUser(to) {
		return errors.New("target user not exists"), ""
	}
	err := t.svc.SendMsg(t.svc.GetUser(to), msg, t.CurrentUser.Name)
	if err != nil {
		return err, ""
	}
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
		t.Context = &TenantContext{
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

func (t *Tenant) Bye() {
	if t.CurrentUser != nil {
		_ = t.svc.UnRegisterTenant(t.CurrentUser.Name)
	}
}
