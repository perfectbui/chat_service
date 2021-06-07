package enum

type ActionValue string

type EnumActionType struct {
	SEND_MESSAGE ActionValue
	JOIN_ROOM    ActionValue
	LEAVE_ROME   ActionValue
}

var Action = EnumActionType{"SEND_MESSAGE", "JOIN_ROOM", "LEAVE_ROOM"}
