package room

import (
	//"log"
)

type Room interface {
	Init(string)
	Close()
	GetType() string
	AddClient() Client
	DelClient(int)
}

type Client interface {
	GetRoom() Room
	Close()
	ReadString(string)
	SetFun_WriteString(func(string))
	SetFun_Disconnect(func())
}

type MatcherParams struct {
}

var (
	roomMap = make(map[string]Room)
)

func DeleteRoom(roomId string) {
	delete(roomMap, roomId)
}
 
// 根据type创建不同的房间
func createRoom(roomId, roomType string) Room {
	var room Room
	switch {
		case roomType == "1" :
			room = new(Room1)
			room.Init(roomId)
		default:
			return nil
	}
	return room
}

// 匹配房间
func Matcher(roomId, roomType string,  params *MatcherParams) Client {
	var room Room
	room = roomMap[roomId]
	if room != nil {
		if room.GetType() != roomType { return nil }
	} else {
			room = createRoom(roomId, roomType)
	}

	if room == nil {
		return nil
	}
	
	client := room.AddClient()
	roomMap[roomId] = room
	return client
}