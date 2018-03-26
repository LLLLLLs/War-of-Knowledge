package internal

var (
	rooms      = make(map[int]*Room)
	lastRoomId = int(0)
)

func AddRoom(room *Room) {
	rooms[room.RoomId] = room
}

func DeleteRoom(roomId int) {
	delete(rooms, roomId)
}

func GetRoom(roomId int) *Room {
	return rooms[roomId]
}
