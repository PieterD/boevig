package deebee

type DB struct {
	files []*dbFile
}

type dbFile struct {
	id   uint64
	size uint64
}
