package internal

var (
	Rooms      = make(map[int]*Room)
	LastRoomId = int(0)
)

func AddRoom(room *Room) {
	Rooms[room.RoomId] = room
}

func DeleteRoom(roomId int) {
	delete(Rooms, roomId)
}

func GetRoom(roomId int) (*Room, bool) {
	room, ok := Rooms[roomId]
	return room, ok
}
