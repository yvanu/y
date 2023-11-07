package y

func isBeat(typ uint32) bool {
	return typ == 0x80000001 || typ == UndoReply(typ)
}

func isRegister(typ uint32) bool {
	return typ == 0x80000002 || typ == UndoReply(typ)
}

func IsNeedReply(typ uint32) bool {
	return typ&0x80000000 == 0x80000000
}

func AddReply(typ uint32) uint32 {
	return typ & 0xffffffff
}

func UndoReply(typ uint32) uint32 {
	return typ & 0x7fffffff
}

var Id = 0

func genId() uint32 {
	Id += 1
	return uint32(Id)
}

func GenId() uint32 {
	return genId()
}
