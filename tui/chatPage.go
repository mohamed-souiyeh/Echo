package tui

import (
	"context"
	db "echo/db/repository"
	"echo/tui/commands"
	"echo/tui/keymaps"
	"echo/tui/messages"
	"echo/tui/styles"
	"echo/utils"
	"echo/workers"
	"echo/workers/comunication"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type ChatModel struct {
	width  int
	height int

	hubReqChan   chan workers.ClientHubReq
	listningChan chan comunication.RoomMessage

	rooms map[workers.RoomID]chan<- comunication.ClientMessage

	msgs        []string
	CurrentUser db.User

	canSend bool

	generalRessourcesCtx    context.Context
	generalRessourcesCancel context.CancelFunc

	msgsViewPort viewport.Model
	gap          string
	textArea     textarea.Model
	help         help.Model
}

func InitChatModel(hubReqChan chan workers.ClientHubReq, win Window) ChatModel {
	ta := textarea.New()

	ta.Prompt = "┃ "
	ta.Placeholder = keymaps.ChatKeyMaps.Submit.Help().Key + " to send a message..."
	ta.ShowLineNumbers = false
	ta.CharLimit = 42
	ta.MaxHeight = 2
	ta.SetHeight(2)
	ta.SetWidth(win.Width)
	ta.FocusedStyle.Prompt = styles.FocusedStyle
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.Prompt = styles.BlurredStyle
	ta.BlurredStyle.CursorLine = styles.BlurredStyle
	ta.Focus()

	vp := viewport.New(win.Width, win.Height)

	vp.MouseWheelEnabled = true
	vp.SetContent("Welcome to a new chat or the continuation of a forgotten one\n")

	ctx, cancel := context.WithCancel(context.Background())

	return ChatModel{
		width:                   win.Width,
		height:                  win.Height,
		hubReqChan:              hubReqChan,
		listningChan:            make(chan comunication.RoomMessage, 4096),
		rooms:                   make(map[workers.RoomID]chan<- comunication.ClientMessage),
		msgs:                    make([]string, 4096),
		CurrentUser:             db.User{},
		canSend:                 true,
		generalRessourcesCtx:    ctx,
		generalRessourcesCancel: cancel,
		msgsViewPort:            vp,
		gap:                     "\n",
		textArea:                ta,
		help:                    help.New(),
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	if msg, ok := msg.(tea.KeyMsg); ok && !key.Matches(msg, keymaps.ChatKeyMaps.Submit) {
		m.textArea, tiCmd = m.textArea.Update(msg)
	}

	if msg, ok := msg.(tea.KeyMsg); !m.textArea.Focused() || !ok {
		m.msgsViewPort, vpCmd = m.msgsViewPort.Update(msg)
	}

	switch msg := msg.(type) {
	// Handle Window Size Changes
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.help.Width = msg.Width

		m.msgsViewPort.Width = msg.Width
		m.msgsViewPort.Height = msg.Height - m.textArea.Height() - lipgloss.Height(m.gap)
		if len(m.msgs) > 0 {
			// Wrap content before setting it.
			m.msgsViewPort.SetContent(styles.ClientRenderer.NewStyle().Width(m.msgsViewPort.Width).Render(strings.Join(m.msgs, "\n")))
		}
		m.msgsViewPort.GotoBottom()

		m.textArea.SetWidth(msg.Width)
	case messages.LogoutMsg:
		m.generalRessourcesCancel()
		m.listningChan = make(chan comunication.RoomMessage, 4096)
		m.rooms = make(map[workers.RoomID]chan<- comunication.ClientMessage)
		m.msgs = make([]string, 4096)
		m.generalRessourcesCtx, m.generalRessourcesCancel = context.WithCancel(context.Background())
		m.textArea.Reset()
		m.msgsViewPort.SetContent("")
		return m, nil
	case messages.AccessChatMsg:
		m.CurrentUser = msg.User
		// i know i dont need to do this assignment since i am not using nameless functions,
		// but just in case future me try to test something involving closures (i know that will happen), it wont be a problem, .
		listningChan := m.listningChan
		currentUserID := m.CurrentUser.ID
		// still keeping the central hub architectuer because i may add DMs in the future, or even groups maybe
		req := workers.ClientHubReq{
			ClientId:       currentUserID,
			ReqType:        workers.JOIN,
			RoomId:         utils.LobbyRoomId,
			ClientChan:     listningChan,
			Msg:            "lemme iiiiin",
			ResponseChanne: make(chan workers.ClientHubRes, 1),
		}
		cmd := tea.Batch(commands.ChatMsgsListenerCmd(m.generalRessourcesCtx, m.listningChan, m.CurrentUser.ID),
			commands.JoinLobbyChatRoomCmd(m.hubReqChan, req))
		return m, tea.Batch(tiCmd, vpCmd, textarea.Blink, cmd)
	case messages.JoinedChatRoomsMsg:
		m.rooms[msg.RoomId] = msg.RoomChan
		// m.joinMsg = fmt.Sprintf("%s with ID: %d", msg.Msg, msg.RoomId)
		return m, nil
	case messages.AllowMsgSend:
		m.canSend = !m.canSend
		return m, nil
	case messages.ChatMsg:
		m.msgs = append(m.msgs, styles.SenderStyle.Render(msg.SenderUsername+": ")+msg.Msg)
		m.msgsViewPort.SetContent(styles.ClientRenderer.NewStyle().Width(m.msgsViewPort.Width).Render(strings.Join(m.msgs, "\n")))
		m.msgsViewPort.GotoBottom()
		return m, tea.Batch(tiCmd, vpCmd, commands.ChatMsgsListenerCmd(m.generalRessourcesCtx, m.listningChan, m.CurrentUser.ID))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymaps.ChatKeyMaps.Logout):
			return m, tea.Batch(tiCmd, vpCmd, commands.LogoutCmd(m.CurrentUser.Username))
		case key.Matches(msg, keymaps.ChatKeyMaps.Esc):
			var cmd tea.Cmd = nil
			if m.textArea.Focused() {
				m.textArea.Blur()
			} else {
				cmd = m.textArea.Focus()
			}
			return m, tea.Batch(tiCmd, vpCmd, cmd)
		case key.Matches(msg, keymaps.ChatKeyMaps.Help):
			if !m.textArea.Focused() {
				m.help.ShowAll = !m.help.ShowAll
			}
		case key.Matches(msg, keymaps.ChatKeyMaps.Submit):
			textMsg := strings.TrimSpace(m.textArea.Value())
			if !m.canSend || textMsg == ""{
				return m, tea.Batch(tiCmd, vpCmd)
			}
			clientID := m.CurrentUser.ID
			clientUsername := m.CurrentUser.Username
			sendMsgCmd := func() tea.Msg {
				msg := comunication.ClientMessage{
					ClientId:       clientID,
					ClientUsername: clientUsername,
					Msg:            textMsg,
				}

				m.rooms[utils.LobbyRoomId] <- msg

				// TODO i need to launch a msg that starts the msg percistence cmd
				return nil
			}
			m.textArea.Reset()
			m.canSend = !m.canSend
			return m, tea.Batch(tiCmd, vpCmd, sendMsgCmd, commands.SetTimeoutCmd(msgTimeoutDuration, messages.AllowMsgSend{}))
		default:
			log.Debugf("the key is : %v", msg)
		}
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

/*
NOTE:
in the table that gonna save the msgs we need a colum that gonna have the msg in the print format
which is coloured sender user name and any formating needed, this gonna make formating the msg only once
and printing huge chunk of msg easer (well i hope, we cant know for sure untill benchmarking),
but i am afraid from one thing which the formating being a bit finiky from a terminal to another
which gonna make things not stable.

*/

func (m ChatModel) View() string {
	help := m.help.View(keymaps.ChatKeyMaps)

	body := lipgloss.JoinVertical(lipgloss.Center, m.msgsViewPort.View(), m.gap, m.textArea.View())

	return lipgloss.JoinVertical(
		lipgloss.Center,
		fmt.Sprintf("chat current user is => %s", m.CurrentUser.Username),
		body,
		help,
	)
}
