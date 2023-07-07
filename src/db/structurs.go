package db

import (
	"strings"
	"time"
)

type User struct {
	ID          int64     `db:"tid"`
	Name        string    `db:"name"`
	UserName    string    `db:"user_name"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Context     int       `db:"context"`
	MessageType int       `db:"expected_message_type"`
	Date        time.Time `db:"date_of_registration"`
}

type Context struct {
	Id     int       `db:"id"`
	Name   string    `db:"name"`
	Text   string    `db:"text"`
	Select bool      `db:"selected"`
	Date   time.Time `db:"dateOfCreation"`
}

type Message struct {
	Text string    `db:"text"`
	Date time.Time `db:"dateOfCreation"`
}

type WaitingUser struct {
	Id         int    `db:"id"`
	ByPassword bool   `db:"by_password"`
	UserName   string `db:"user_name"`
	Key        string `db:"key"`
}

func (u *User) TrimSpace() {
	u.Name = strings.TrimSpace(u.Name)
	u.UserName = strings.TrimSpace(u.UserName)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
}

func (u *User) GetDateStr() string {
	return u.Date.Format("02.01.2006 15:04")
}

func (c *Context) TrimSpace() {
	c.Name = strings.TrimSpace(c.Name)
	c.Text = strings.TrimSpace(c.Text)
}

func (c *Context) GetDateStr() string {
	return c.Date.Format("02.01.2006 15:04")
}

func (m *Message) TrimSpace() {
	m.Text = strings.TrimSpace(m.Text)
}

func (m *Message) GetDateStr() string {
	return m.Date.Format("02.01.2006 15:04")
}

func (w *WaitingUser) TrimSpace() {
	w.UserName = strings.TrimSpace(w.UserName)
	w.Key = strings.TrimSpace(w.Key)
}
