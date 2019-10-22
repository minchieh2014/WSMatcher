package room

import (
	"log"
	"time"
)

type Client1 struct{
	Id int
	Room *Room1
	WriteString func(string)
	Disconnect func()
}

type Room1 struct {
	Id string
	clientMap map[int]*Client1 
}

const MaxClient = 2
var clientIdInc int

/***************Room*******************/
func (r *Room1)GetType() string { return "1" }
func (r *Room1)AddClient() Client {
	if len(r.clientMap) >= MaxClient { return nil }
	c := new(Client1)
	clientIdInc++
	c.Id = clientIdInc
	c.Room = r
	r.clientMap[c.Id] = c
	log.Println("Room1:", r.Id, "new client", c.Id, "count", len(r.clientMap))
	if len(r.clientMap) == 1 { // room1 long time only one client -> close room
		go func() {
			time.Sleep(time.Duration(10)*time.Second)
			if len(r.clientMap) != 2 {
				r.Close()
			}
		}()
	}
	return c
}

func (r *Room1) DelClient(id int) {
		// room1: del one -> close room
		r.Close()
}

func (r *Room1) Init(id string) {
	log.Println("Room1:", id, "new room")
	r.Id = id
	r.clientMap = make(map[int]*Client1)
}

func (r *Room1) Close() {
	log.Println("Room1:", r.Id, "close")
	for id, c := range r.clientMap {
		if c != nil {
			c.Room = nil
			c.Close()
		}
		delete(r.clientMap,id)
		log.Println("Room1:", r.Id, "del client", id, "count", len(r.clientMap))
	}
	DeleteRoom(r.Id)
}

/****************Client******************/

func (c *Client1)GetRoom() Room {
	return c.Room
}

func (c *Client1)Close() {
	if c.Room != nil {
		r := c.Room
		c.Room = nil
		r.DelClient(c.Id)
	}
	if c.Disconnect != nil { c.Disconnect() }
}

func (c *Client1)ReadString(str string) {
	r := c.GetRoom().(*Room1)
	for id,cIter := range r.clientMap {
		if id != c.Id {
			cIter.WriteString(str)
		}
	}
}

func (c *Client1)SetFun_WriteString(f func(string)) {
	c.WriteString = f
}

func (c *Client1)SetFun_Disconnect(f func()) {
	c.Disconnect = f
}
