package aspect

type Feature struct {
	Name        string
	Passable    bool
	Transparent bool
	Wallable    bool
}

var Void = Feature{}
var Wall = Feature{
	Name:     "wall",
	Wallable: true,
}
var Floor = Feature{
	Name:        "floor",
	Passable:    true,
	Transparent: true,
}
var ClosedDoor = Feature{
	Name:     "closed door",
	Wallable: true,
}
var OpenDoor = Feature{
	Name:        "open door",
	Wallable:    true,
	Passable:    true,
	Transparent: true,
}
