package chat

type User interface {
	GetId() int64
	SendMsg(msg IMsg) error
}

type ChatRoom struct {
	Id            string
	users         map[int64]User
	join          chan User
	leave         chan User
	broadcast     chan IMsg
	joinCallback  func(user User)
	leaveCallback func(user User)
	processMsg    func(msg IMsg)
}

func (c *ChatRoom) OnJoin(f func(user User)) {
	c.joinCallback = f
}

func (c *ChatRoom) OnLeave(f func(user User)) {
	c.leaveCallback = f
}

func (c *ChatRoom) ProcessMsg(f func(msg IMsg)) {
	c.processMsg = f
}

func (c *ChatRoom) Join(user User) {
	c.join <- user
}

func (c *ChatRoom) Leave(user User) {
	c.leave <- user
}

func (c *ChatRoom) Broadcast(m IMsg) {
	c.broadcast <- m
}

func (c *ChatRoom) Run() {
	for {
		select {
		case u := <-c.join:
			c.users[u.GetId()] = u
			if c.joinCallback != nil {
				c.joinCallback(u)
			}
		case u := <-c.leave:
			delete(c.users, u.GetId())
			if c.leaveCallback != nil {
				c.leaveCallback(u)
			}
		case m := <-c.broadcast:
			c.processMsg(m)
			for _, u := range c.users {
				//if u.GetId() != m.GetUser().GetId() { // 排除自己
				u.SendMsg(m)
				//}
			}
		}
	}
}

func (c *ChatRoom) UserCount() int {
	return len(c.users)
}

func NewChatRoom(Id string) *ChatRoom {
	return &ChatRoom{
		Id:        Id,
		users:     make(map[int64]User),
		join:      make(chan User),
		leave:     make(chan User),
		broadcast: make(chan IMsg),
	}
}
