package frontend

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type OperationType int

const (
	SEND_MESSAGE OperationType = iota
	CREATE_CHAT
	JOIN_CHAT
	LEAVE_CHAT
	INVITE_TO_CHAT
	SEND_FILE
	SET_USERNAME
	TEST_MESSAGE
	LOAD_MESSAGES
)

type FrontendMessage struct {
	Timestamp int64
	Content   string // Message content
	FromUser  string // UserID
	ChatID    string // ChatID
	Operation OperationType
}

type screen int

const (
	screenIntro screen = iota
	screenChats
	screenHelp
	screenUsername
)

type model struct {
	currentScreen screen                       // Current screen
	chats         map[string][]FrontendMessage // Map from ChatID to messages
	chatNames     map[string]string            // Map from ChatID to chat names
	invites       []string                     // List of invites
	cursor        int                          // Cursor for selecting chats or invites
	focus         string                       // Can be 'chats' or 'invites'
	currentChat   string                       // Currently selected chat
	inChatDetail  bool                         // Whether we are in chat details
	input         textinput.Model              // User input for messages
	usernameInput textinput.Model              // User input for setting the username
	usernames     map[string]string            // Map from ChatID to username
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 20

	ui := textinput.New()
	ui.Placeholder = "Set your username..."
	ui.CharLimit = 20
	ui.Width = 20

	return model{
		currentScreen: screenIntro,
		chats: map[string][]FrontendMessage{
			"1": {
				{Timestamp: time.Now().Unix(), Content: "Hello, this is chat 1!", FromUser: "Alice", ChatID: "1"},
				{Timestamp: time.Now().Unix(), Content: "Hi Alice!", FromUser: "Bob", ChatID: "1"},
			},
			"2": {
				{Timestamp: time.Now().Unix(), Content: "Welcome to chat 2", FromUser: "Charlie", ChatID: "2"},
			},
		},
		chatNames:     map[string]string{"1": "Chat 1", "2": "Chat 2"},
		invites:       []string{"Invitation from Alice to Chat 3", "Invitation from Bob to Chat 4"},
		focus:         "chats",
		input:         ti,
		usernameInput: ui,
		usernames:     make(map[string]string),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func createMessage(fromUser, chatID, content string, op OperationType) FrontendMessage {
	return FrontendMessage{
		Timestamp: time.Now().Unix(),
		Content:   content,
		FromUser:  fromUser,
		ChatID:    chatID,
		Operation: op,
	}
}

func isValidUsername(username string) bool {
	if len(username) > 20 {
		return false
	}
	for _, r := range username {
		if unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func (m *model) handleCommand(command, argument string) (tea.Model, tea.Cmd) {
	var msg FrontendMessage

	switch command {
	case "/test":
		msg = createMessage(m.usernames[m.currentChat], m.currentChat, argument, TEST_MESSAGE)
	case "/leave":
		msg = createMessage(m.usernames[m.currentChat], m.currentChat, "", LEAVE_CHAT)
		m.chats[m.currentChat] = append(m.chats[m.currentChat], msg)
		m.inChatDetail = false
		delete(m.chats, m.currentChat)
		delete(m.chatNames, m.currentChat)
		delete(m.usernames, m.currentChat)
		m.currentChat = ""
	case "/invite":
		msg = createMessage(m.usernames[m.currentChat], m.currentChat, argument, INVITE_TO_CHAT)
	case "/sendfile":
		msg = createMessage(m.usernames[m.currentChat], m.currentChat, argument, SEND_FILE)
	case "/setusername":
		if isValidUsername(argument) {
			oldUsername := m.usernames[m.currentChat]
			m.usernames[m.currentChat] = argument
			msg = createMessage(oldUsername, m.currentChat, fmt.Sprintf("%s changed their username to %s", oldUsername, argument), SET_USERNAME)
		} else {
			msg = createMessage(m.usernames[m.currentChat], m.currentChat, "Invalid username: "+argument, SEND_MESSAGE)
		}
	case "/loadmessages":
		msg = createMessage(m.usernames[m.currentChat], m.currentChat, "Loading more messages...", LOAD_MESSAGES)
	default:
		msg = createMessage(m.usernames[m.currentChat], m.currentChat, "Unknown command: "+command, SEND_MESSAGE)
	}

	m.chats[m.currentChat] = append(m.chats[m.currentChat], msg)
	m.input.SetValue("") // Clear the input after handling the command
	return m, nil
}

func (m *model) handleEnter() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.input.Value())
	if len(input) > 0 {
		if strings.HasPrefix(input, "/") {
			parts := strings.SplitN(input, " ", 2)
			command := parts[0]
			var argument string
			if len(parts) > 1 {
				argument = parts[1]
			}
			return m.handleCommand(command, argument)
		} else {
			msg := createMessage(m.usernames[m.currentChat], m.currentChat, input, SEND_MESSAGE)
			m.chats[m.currentChat] = append(m.chats[m.currentChat], msg)
			m.input.SetValue("") // Clear the input after sending a message
		}
	}
	return m, nil
}

func (m *model) handleInviteAccept() (tea.Model, tea.Cmd) {
	invitation := m.invites[m.cursor] // Save the current invitation
	parts := strings.Split(invitation, " to ")
	chatID := parts[len(parts)-1]
	chatName := "Chat " + chatID

	// Add the chat if it doesn't exist
	if _, exists := m.chatNames[chatID]; !exists {
		m.chatNames[chatID] = chatName
		m.chats[chatID] = []FrontendMessage{}
	}

	// Remove the accepted invitation and adjust the cursor
	m.invites = append(m.invites[:m.cursor], m.invites[m.cursor+1:]...)
	if len(m.invites) == 0 {
		m.cursor = 0
		m.focus = "chats" // Switch back to "chats" if there are no more invitations
	} else if m.cursor >= len(m.invites) {
		m.cursor = len(m.invites) - 1 // Adjust the cursor if it goes out of bounds
	}
	return m, tea.Printf("Invitation accepted: %s", invitation)
}

func (m *model) handleInviteDecline() (tea.Model, tea.Cmd) {
	m.invites = append(m.invites[:m.cursor], m.invites[m.cursor+1:]...)
	if len(m.invites) == 0 {
		m.cursor = 0
		m.focus = "chats" // Switch back to "chats" if there are no more invitations
	} else if m.cursor >= len(m.invites) {
		m.cursor = len(m.invites) - 1 // Adjust the cursor if it goes out of bounds
	}
	return m, tea.Printf("Invitation declined")
}

func (m *model) handleChatSelection() (tea.Model, tea.Cmd) {
	chatIDs := getSortedChatIDs(m.chatNames)
	if m.cursor < len(chatIDs) {
		m.currentChat = chatIDs[m.cursor]
		if _, exists := m.usernames[m.currentChat]; !exists {
			m.currentScreen = screenUsername
			m.usernameInput.Focus()
		} else {
			m.inChatDetail = true
			m.input.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentScreen {
		case screenIntro:
			if msg.String() == "enter" || msg.String() == " " {
				m.currentScreen = screenChats
			}
		case screenChats:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				if m.inChatDetail {
					m.inChatDetail = false
					m.input.Blur()
				} else {
					return m, tea.Quit
				}
			case tea.KeyRunes:
				if msg.String() == "q" && !m.inChatDetail {
					return m, tea.Quit
				}
				if msg.String() == "h" && !m.inChatDetail {
					m.currentScreen = screenHelp
				}
			case tea.KeyUp:
				if !m.inChatDetail && m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyDown:
				if !m.inChatDetail {
					if m.focus == "chats" && m.cursor < len(m.chatNames)-1 {
						m.cursor++
					} else if m.focus == "invites" && m.cursor < len(m.invites)-1 {
						m.cursor++
					}
				}
			case tea.KeyTab:
				if !m.inChatDetail {
					m.focus = toggleFocus(m.focus)
					if len(m.invites) > 0 {
						m.cursor = 0
					}
				}
			case tea.KeyEnter:
				if m.inChatDetail {
					return m.handleEnter()
				} else if m.focus == "invites" && len(m.invites) > 0 {
					return m.handleInviteAccept()
				} else if m.focus == "chats" {
					return m.handleChatSelection()
				}
			case tea.KeyBackspace:
				if m.focus == "invites" && len(m.invites) > 0 {
					model, cmd := m.handleInviteDecline()
					return model, cmd
				}
			}

			if m.inChatDetail {
				return m.handleChatInput(msg)
			}
		case screenHelp:
			if msg.Type == tea.KeyEsc {
				m.currentScreen = screenChats
			}
		case screenUsername:
			switch msg.Type {
			case tea.KeyEnter:
				username := strings.TrimSpace(m.usernameInput.Value())
				if isValidUsername(username) {
					m.usernames[m.currentChat] = username
					m.usernameInput.Blur()
					m.currentScreen = screenChats
					m.inChatDetail = true
					m.input.Focus()
					joinMsg := createMessage(username, m.currentChat, fmt.Sprintf("%s joined the chat", username), JOIN_CHAT)
					m.chats[m.currentChat] = append(m.chats[m.currentChat], joinMsg)
					m.usernameInput.SetValue("")
					return m, textinput.Blink
				} else {
					return m, tea.Printf("Invalid username. Please enter a single word with up to 20 characters.")
				}
			}

			var cmd tea.Cmd
			m.usernameInput, cmd = m.usernameInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) handleChatInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyLeft && len(m.input.Value()) == 0 {
		return m, nil // Ignore left arrow key if the input field is empty
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func toggleFocus(focus string) string {
	if focus == "chats" {
		return "invites"
	}
	return "chats"
}

func (m model) View() string {
	switch m.currentScreen {
	case screenIntro:
		return m.introView()
	case screenChats:
		return m.chatsView()
	case screenHelp:
		return m.helpView()
	case screenUsername:
		return m.usernameView()
	}
	return ""
}

func (m model) introView() string {
	return `
.▄▄ · ▄ •▄ ▄• ▄▌ ▐ ▄ ▄ •▄ 
▐█ ▀. █▌▄▌▪█▪██▌•█▌▐██▌▄▌▪     ^...^
▄▀▀▀█▄▐▀▀▄·█▌▐█▌▐█▐▐▌▐▀▀▄·    <_* *_>   
▐█▄▪▐█▐█.█▌▐█▄█▌██▐█▌▐█.█▌      \_/
 ▀▀▀▀ ·▀  ▀ ▀▀▀ ▀▀ █▪·▀  ▀

Welcome to Skunk, a cutting-edge peer-to-peer communication platform that enables secure and private messaging with powerful encryption.

Press Enter to continue...
`
}

func (m model) chatsView() string {
	if m.inChatDetail {
		return m.chatDetailView()
	}

	s := `
╔════════════════════════════════════════════════════════════╗
║                           Skunk                            ║
║       Secure Peer-to-Peer Communication Platform           ║
╚════════════════════════════════════════════════════════════╝
`

	s += "╔════════════════════════════════════════════════════════════╗\n"
	s += "║            Chats                                           ║\n"
	s += "╠════════════════════════════════════════════════════════════╣\n"

	chatIDs := getSortedChatIDs(m.chatNames)

	for i, chatID := range chatIDs {
		cursor := " " // Space if the cursor is not on this line
		if i == m.cursor && m.focus == "chats" {
			cursor = ">" // Show the cursor
		}
		s += fmt.Sprintf("║ %s %s\n", cursor, m.chatNames[chatID])
	}

	s += "╠════════════════════════════════════════════════════════════╣\n"
	s += "║           Invites                                          ║\n"
	s += "╠════════════════════════════════════════════════════════════╣\n"

	for i, invite := range m.invites {
		cursor := " "
		if i == m.cursor && m.focus == "invites" {
			cursor = ">"
		}
		s += fmt.Sprintf("║ %s %s\n", cursor, invite)
	}

	if m.focus == "invites" && len(m.invites) == 0 {
		s += "║ No invitations available                               ║\n"
	}

	s += "╚════════════════════════════════════════════════════════════╝\n"
	s += "\nPress Enter to open the chat, Tab to switch between chats and invites, Backspace to decline an invite, h for help, and q to quit."

	return s
}

func (m model) chatDetailView() string {
	s := fmt.Sprintf("Chat: %s\n\n", m.chatNames[m.currentChat])

	for _, msg := range m.chats[m.currentChat] {
		timeString := time.Unix(msg.Timestamp, 0).Format("2006-01-02 15:04:05")
		switch msg.Operation {
		case SEND_MESSAGE:
			s += fmt.Sprintf("[%s] %s: %s\n", timeString, msg.FromUser, msg.Content)
		case CREATE_CHAT:
			s += fmt.Sprintf("[%s] Chat created by %s with ChatID %s\n", timeString, msg.FromUser, msg.ChatID)
		case JOIN_CHAT:
			s += fmt.Sprintf("[%s] %s\n", timeString, msg.Content)
		case LEAVE_CHAT:
			s += fmt.Sprintf("[%s] User %s has left chat %s\n", timeString, msg.FromUser, msg.ChatID)
		case INVITE_TO_CHAT:
			s += fmt.Sprintf("[%s] User %s has been invited by %s\n", timeString, msg.Content, msg.FromUser)
		case SEND_FILE:
			s += fmt.Sprintf("[%s] User %s has sent a file in chat %s\n", timeString, msg.FromUser, msg.ChatID)
		case SET_USERNAME:
			s += fmt.Sprintf("[%s] %s\n", timeString, msg.Content)
		case TEST_MESSAGE:
			s += fmt.Sprintf("[%s] Test message from %s: %s\n", timeString, msg.FromUser, msg.Content)
		case LOAD_MESSAGES:
			s += fmt.Sprintf("[%s] %s: %s\n", timeString, msg.FromUser, msg.Content)
		default:
			s += fmt.Sprintf("[%s] Unknown operation received from %s\n", timeString, msg.FromUser)
		}
	}

	s += fmt.Sprintf("\n%s", m.input.View())
	s += "\nPress Enter to send the message, ESC to go back."

	return s
}

func (m model) helpView() string {
	return `
Available commands in the chat:

  /test <OnionID> - Send a test message to verify connectivity
  /leave - Leave a chat
  /invite <OnionID> - Invite a user to a chat
  /sendfile <FilePath> - Send a file in a chat
  /setusername <NewUsername> - Set or change the user's username
  /loadmessages - Loads 50 more messages of the chat if they exist

Press ESC to return to the main menu.
`
}

func (m model) usernameView() string {
	return fmt.Sprintf("Please set your username for %s:\n\n%s\n\nPress Enter to confirm.", m.chatNames[m.currentChat], m.usernameInput.View())
}

func getSortedChatIDs(chatNames map[string]string) []string {
	chatIDs := make([]string, 0, len(chatNames))
	for chatID := range chatNames {
		chatIDs = append(chatIDs, chatID)
	}
	sort.Strings(chatIDs) // Sort chat IDs to ensure a stable order
	return chatIDs
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Oh no, something went wrong: %s\n", err)
		os.Exit(1)
	}
}
