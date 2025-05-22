package commands

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)


func SetTimeoutCmd(time time.Duration, msg interface{}) tea.Cmd{
	ctx, _ := context.WithTimeout(context.Background(), time)

	return func() tea.Msg {

		<- ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			return msg
		}
		return nil
	}
}