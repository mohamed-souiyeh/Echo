package tui

// import "fmt"

/*
	experemant with the charmbracelet logger to achive what u did here i would make a nice
	little library to use in the future
*/

type debugger struct {
	Enabled bool
}

var debug debugger = debugger{
	Enabled: false,
}


func (d debugger) print(call func ()) {
	if d.Enabled {
		call()
	}
}