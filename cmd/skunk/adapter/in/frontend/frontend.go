package frontend

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// OperationType represents different types of operations that can be performed.
type OperationType int

// Define various operation types.
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

// FrontendMessage represents a message in the frontend.
type FrontendMessage struct {
	Id              string
	Timestamp       int64
	Content         string
	SenderID        string
	ReceiverID      string
	SenderAddress   string
	ReceiverAddress string
	ChatID          string
	Operation       OperationType
}

// Define various screens in the application.
type screen int

const (
	screenIntro screen = iota
	screenChats
	screenHelp
	screenUsername
	screenTestUser
	screenCreateChat
)

// TempMsgTimeoutMsg represents a timeout message for temporary messages.
type TempMsgTimeoutMsg struct{}

// Model represents the application's state.
type Model struct {
	currentScreen     screen                       // Current screen
	Chats             map[string][]FrontendMessage // Map from ChatID to messages
	chatNames         map[string]string            // Map from ChatID to chat names
	invites           []FrontendMessage            // List of invites
	cursor            int                          // Cursor for selecting Chats or invites
	focus             string                       // Can be 'Chats' or 'invites'
	CurrentChat       string                       // Currently selected chat
	inChatDetail      bool                         // Whether we are in chat details
	input             textinput.Model              // User input for messages
	usernameInput     textinput.Model              // User input for setting the username
	Usernames         map[string]string            // Map from ChatID to username
	testUserInput     textinput.Model              // User input for testing connection
	createChatInput   textinput.Model              // User input for creating chat (invitees)
	chatNameInput     textinput.Model              // User input for creating chat (chat name)
	chatInvitees      []string                     // List of invitees for the new chat
	TempMessage       string                       // Temporary message
	TempMessageExpire time.Time                    // Expiry time for the temporary message
}

// InitialModel returns the initial state of the application.
func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40

	ui := textinput.New()
	ui.Placeholder = "Set your username..."
	ui.CharLimit = 20
	ui.Width = 40

	tu := textinput.New()
	tu.Placeholder = "Enter OnionID to test connection..."
	tu.CharLimit = 56
	tu.Width = 40

	ci := textinput.New()
	ci.Placeholder = "Enter OnionID to invite..."
	ci.CharLimit = 56
	ci.Width = 100

	cn := textinput.New()
	cn.Placeholder = "Enter chat name..."
	cn.CharLimit = 30
	cn.Width = 40

	return Model{
		currentScreen:   screenIntro,
		Chats:           map[string][]FrontendMessage{},
		chatNames:       map[string]string{},
		invites:         []FrontendMessage{},
		focus:           "Chats",
		input:           ti,
		usernameInput:   ui,
		Usernames:       make(map[string]string),
		testUserInput:   tu,
		createChatInput: ci,
		chatNameInput:   cn,
		chatInvitees:    []string{},
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// CreateMessage creates a new FrontendMessage.
func CreateMessage(chatID, senderID, receiverID, content, senderAddress, receiverAddress string, op OperationType) FrontendMessage {
	return FrontendMessage{
		Id:              fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp:       time.Now().Unix(),
		Content:         SanitizeInput(content),
		SenderID:        SanitizeInput(senderID),
		ReceiverID:      SanitizeInput(receiverID),
		SenderAddress:   SanitizeInput(senderAddress),
		ReceiverAddress: SanitizeInput(receiverAddress),
		ChatID:          chatID,
		Operation:       op,
	}
}

// IsValidUsername checks if a username is valid.
func IsValidUsername(username string) bool {
	if len(username) == 0 || len(username) > 20 {
		return false
	}
	for _, r := range username {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) || unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

// SanitizeInput sanitizes a string input.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// ValidateInput validates and sanitizes a string input.
func ValidateInput(input string, maxLength int) (string, error) {
	sanitized := SanitizeInput(input)
	if len(sanitized) > maxLength {
		return "", errors.New("input exceeds maximum length")
	}
	for _, r := range sanitized {
		if !unicode.IsPrint(r) {
			return "", errors.New("input contains non-printable characters")
		}
	}
	return sanitized, nil
}

// HandleCommand handles user commands.
func (m *Model) HandleCommand(command, argument string) (tea.Model, tea.Cmd) {
	var msg FrontendMessage
	var err error

	switch command {
	case "/leave":
		msg = CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], "", "", "", "", LEAVE_CHAT)
		m.Chats[m.CurrentChat] = append(m.Chats[m.CurrentChat], msg)
		m.inChatDetail = false
		delete(m.Chats, m.CurrentChat)
		delete(m.chatNames, m.CurrentChat)
		delete(m.Usernames, m.CurrentChat)
		m.CurrentChat = ""
		return m, nil
	case "/invite":
		argument, err = ValidateInput(argument, 56)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		msg = CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], argument, argument, "", "", INVITE_TO_CHAT)
		m.invites = append(m.invites, msg)
		m.Chats[m.CurrentChat] = append(m.Chats[m.CurrentChat], msg)
		return m, tea.Printf("Invited %s to the chat", argument)
	case "/sendfile":
		argument, err = ValidateInput(argument, 256)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		msg = CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], "", argument, "", "", SEND_FILE)
	case "/setusername":
		argument, err = ValidateInput(argument, 20)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		if IsValidUsername(argument) {
			oldUsername := m.Usernames[m.CurrentChat]
			m.Usernames[m.CurrentChat] = argument
			msg = CreateMessage(m.CurrentChat, oldUsername, argument, fmt.Sprintf("%s changed their username to %s", oldUsername, argument), "", "", SET_USERNAME)
		} else {
			msg = CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], "", "Invalid username: "+argument, "", "", SEND_MESSAGE)
		}
	case "/loadmessages":
		msg = CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], "", "Loading more messages...", "", "", LOAD_MESSAGES)
	default:
		argument, err = ValidateInput(argument, 256)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		msg = CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], "", "Unknown command: "+command, "", "", SEND_MESSAGE)
	}

	m.Chats[m.CurrentChat] = append(m.Chats[m.CurrentChat], msg)
	m.input.SetValue("")
	return m, nil
}

// HandleEnter handles the Enter key press.
func (m *Model) HandleEnter() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.input.Value())
	if len(input) > 0 {
		if strings.HasPrefix(input, "/") {
			parts := strings.SplitN(input, " ", 2)
			command := parts[0]
			var argument string
			if len(parts) > 1 {
				argument = parts[1]
			}
			return m.HandleCommand(command, argument)
		} else {
			msg := CreateMessage(m.CurrentChat, m.Usernames[m.CurrentChat], "", input, "", "", SEND_MESSAGE)
			m.Chats[m.CurrentChat] = append(m.Chats[m.CurrentChat], msg)
			m.input.SetValue("")
		}
	}
	return m, nil
}

// HandleInviteAccept handles accepting an invite.
func (m *Model) HandleInviteAccept() (tea.Model, tea.Cmd) {
	invitation := m.invites[m.cursor]
	chatID := invitation.ChatID
	chatName := chatID

	if _, exists := m.chatNames[chatID]; !exists {
		m.chatNames[chatID] = chatName
		m.Chats[chatID] = []FrontendMessage{}
	}

	m.invites = append(m.invites[:m.cursor], m.invites[m.cursor+1:]...)
	if len(m.invites) == 0 {
		m.cursor = 0
		m.focus = "Chats"
	} else if m.cursor >= len(m.invites) {
		m.cursor = len(m.invites) - 1
	}
	m.TempMessage = fmt.Sprintf("Invitation accepted: %s", invitation.Content)
	m.TempMessageExpire = time.Now().Add(10 * time.Second)
	return m, m.ClearTempMessage()
}

// HandleInviteDecline handles declining an invite.
func (m *Model) HandleInviteDecline() (tea.Model, tea.Cmd) {
	m.invites = append(m.invites[:m.cursor], m.invites[m.cursor+1:]...)
	if len(m.invites) == 0 {
		m.cursor = 0
		m.focus = "Chats"
	} else if m.cursor >= len(m.invites) {
		m.cursor = len(m.invites) - 1
	}
	m.TempMessage = "Invitation declined"
	m.TempMessageExpire = time.Now().Add(10 * time.Second)
	return m, m.ClearTempMessage()
}

// HandleChatSection handles selecting a chat.
func (m *Model) HandleChatSection() (tea.Model, tea.Cmd) {
	chatIDs := GetSortedChatIDs(m.chatNames)
	if m.cursor < len(chatIDs) {
		m.CurrentChat = chatIDs[m.cursor]
		if _, exists := m.Usernames[m.CurrentChat]; !exists {
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

// Update updates the application's state based on incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if msg.String() == "c" && !m.inChatDetail {
					m.currentScreen = screenCreateChat
					m.createChatInput.Focus()
				}
				if msg.String() == "t" && !m.inChatDetail {
					m.currentScreen = screenTestUser
					m.testUserInput.Focus()
				}
			case tea.KeyUp:
				if !m.inChatDetail && m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyDown:
				if !m.inChatDetail {
					if m.focus == "Chats" && m.cursor < len(m.chatNames)-1 {
						m.cursor++
					} else if m.focus == "invites" && m.cursor < len(m.invites)-1 {
						m.cursor++
					}
				}
			case tea.KeyTab:
				if !m.inChatDetail {
					m.focus = ToggleFocus(m.focus)
					if len(m.invites) > 0 {
						m.cursor = 0
					}
				}
			case tea.KeyEnter:
				if m.inChatDetail {
					return m.HandleEnter()
				} else if m.focus == "invites" && len(m.invites) > 0 {
					return m.HandleInviteAccept()
				} else if m.focus == "Chats" {
					return m.HandleChatSection()
				}
			case tea.KeyBackspace:
				if m.focus == "invites" && len(m.invites) > 0 {
					Model, cmd := m.HandleInviteDecline()
					return Model, cmd
				}
			}

			if m.inChatDetail {
				return m.HandleChatInput(msg)
			}
		case screenHelp:
			if msg.Type == tea.KeyEsc {
				m.currentScreen = screenChats
			}
		case screenUsername:
			switch msg.Type {
			case tea.KeyEnter:
				username := strings.TrimSpace(m.usernameInput.Value())
				if IsValidUsername(username) {
					m.Usernames[m.CurrentChat] = username
					m.usernameInput.Blur()
					m.currentScreen = screenChats
					m.inChatDetail = true
					m.input.Focus()
					joinMsg := CreateMessage(m.CurrentChat, username, "", fmt.Sprintf("%s joined the chat", username), "", "", JOIN_CHAT)
					m.Chats[m.CurrentChat] = append(m.Chats[m.CurrentChat], joinMsg)
					m.usernameInput.SetValue("")
					return m, textinput.Blink
				} else {
					return m, tea.Printf("Invalid username. Please enter a single word with up to 20 characters.")
				}
			}

			var cmd tea.Cmd
			m.usernameInput, cmd = m.usernameInput.Update(msg)
			cmds = append(cmds, cmd)
		case screenTestUser:
			switch msg.Type {
			case tea.KeyEnter:
				onionID := strings.TrimSpace(m.testUserInput.Value())
				m.testUserInput.SetValue("")
				if TestConnection(onionID) {
					m.TempMessage = fmt.Sprintf("Connection to %s successful!", onionID)
				} else {
					m.TempMessage = fmt.Sprintf("Failed to connect to %s.", onionID)
				}
				m.TempMessageExpire = time.Now().Add(10 * time.Second)
				return m, m.ClearTempMessage()
			case tea.KeyEsc:
				m.testUserInput.Blur()
				m.currentScreen = screenChats
			}

			var cmd tea.Cmd
			m.testUserInput, cmd = m.testUserInput.Update(msg)
			cmds = append(cmds, cmd)
		case screenCreateChat:
			switch msg.Type {
			case tea.KeyCtrlC:
				chatName := strings.TrimSpace(m.chatNameInput.Value())
				if chatName == "" {
					return m, tea.Printf("Chat name cannot be empty.")
				}
				chatID := fmt.Sprintf("%d", time.Now().UnixNano())
				m.chatNames[chatID] = chatName
				m.Chats[chatID] = []FrontendMessage{}
				for _, invitee := range m.chatInvitees {
					msg := CreateMessage(chatID, m.Usernames[m.CurrentChat], invitee, invitee, "", "", INVITE_TO_CHAT)
					m.Chats[chatID] = append(m.Chats[chatID], msg)
				}
				m.createChatInput.SetValue("")
				m.chatNameInput.SetValue("")
				m.chatInvitees = []string{}
				m.currentScreen = screenChats
				m.focus = "Chats"
				m.TempMessage = fmt.Sprintf("Chat %s created!", chatName)
				m.TempMessageExpire = time.Now().Add(10 * time.Second)
				return m, m.ClearTempMessage()
			case tea.KeyCtrlI:
				if m.createChatInput.Focused() {
					m.createChatInput.Blur()
					m.chatNameInput.Focus()
				} else {
					m.chatNameInput.Blur()
					m.createChatInput.Focus()
				}
			case tea.KeyEnter:
				if m.createChatInput.Focused() {
					invitee := strings.TrimSpace(m.createChatInput.Value())
					if invitee != "" {
						m.chatInvitees = append(m.chatInvitees, invitee)
						m.createChatInput.SetValue("")
					}
				} else if m.chatNameInput.Focused() {
					m.chatNameInput.Blur()
					m.createChatInput.Focus()
				}
			case tea.KeyEsc:
				m.createChatInput.SetValue("")
				m.chatNameInput.SetValue("")
				m.chatInvitees = []string{}
				m.currentScreen = screenChats
			}

			var cmd tea.Cmd
			m.createChatInput, cmd = m.createChatInput.Update(msg)
			cmds = append(cmds, cmd)
			m.chatNameInput, cmd = m.chatNameInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	case TempMsgTimeoutMsg:
		if time.Now().After(m.TempMessageExpire) {
			m.TempMessage = ""
		}
	}

	return m, tea.Batch(cmds...)
}

// HandleChatInput handles chat input.
func (m *Model) HandleChatInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyLeft && len(m.input.Value()) == 0 {
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// ToggleFocus toggles the focus between Chats and invites.
func ToggleFocus(focus string) string {
	if focus == "Chats" {
		return "invites"
	}
	return "Chats"
}

// View returns the view of the application.
func (m Model) View() string {
	tempMsg := ""
	if m.TempMessage != "" && time.Now().Before(m.TempMessageExpire) {
		tempMsg = fmt.Sprintf("\n\n%s\n\n", m.TempMessage)
	}
	switch m.currentScreen {
	case screenIntro:
		return m.IntroView()
	case screenChats:
		return m.ChatsView() + tempMsg
	case screenHelp:
		return m.helpView()
	case screenUsername:
		return m.UsernameView()
	case screenTestUser:
		return m.TestUserView() + tempMsg
	case screenCreateChat:
		return m.CreateChatView()
	}
	return ""
}

// IntroView returns the view for the intro screen.
func (m Model) IntroView() string {
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

// ChatsView returns the view for the Chats screen.
func (m Model) ChatsView() string {
	if m.inChatDetail {
		return m.ChatDetailView()
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

	chatIDs := GetSortedChatIDs(m.chatNames)

	for i, chatID := range chatIDs {
		cursor := " "
		if i == m.cursor && m.focus == "Chats" {
			cursor = ">"
		}
		padding := len(" " + cursor + " " + m.chatNames[chatID])
		spaces := 60 - padding
		s += fmt.Sprintf("║ %s %s%*s║\n", cursor, m.chatNames[chatID], spaces, "")
	}

	s += "╠════════════════════════════════════════════════════════════╣\n"
	s += "║           Invites                                          ║\n"
	s += "╠════════════════════════════════════════════════════════════╣\n"

	for i, invite := range m.invites {
		cursor := " "
		if i == m.cursor && m.focus == "invites" {
			cursor = ">"
		}
		padding := len(" " + cursor + " " + invite.Content)
		spaces := 60 - padding
		s += fmt.Sprintf("║ %s %s%*s║\n", cursor, invite.Content, spaces, "")
	}

	if m.focus == "invites" && len(m.invites) == 0 {
		s += "║ No invitations available                               ║\n"
	}

	s += "╚════════════════════════════════════════════════════════════╝\n"
	s += "\nPress Enter to open the chat, Tab to switch between Chats and invites, Backspace to decline an invite, h for help, t to test connection, c to create a chat, and q to quit."

	return s
}

// ChatDetailView returns the view for the chat detail screen.
func (m Model) ChatDetailView() string {
	s := fmt.Sprintf("Chat: %s\n\n", m.chatNames[m.CurrentChat])

	for _, msg := range m.Chats[m.CurrentChat] {
		timeString := time.Unix(msg.Timestamp, 0).Format("2006-01-02 15:04:05")
		switch msg.Operation {
		case SEND_MESSAGE:
			s += fmt.Sprintf("[%s] %s: %s\n", timeString, msg.SenderID, msg.Content)
		case CREATE_CHAT:
			s += fmt.Sprintf("[%s] Chat created by %s with ChatID %s\n", timeString, msg.SenderID, msg.ChatID)
		case JOIN_CHAT:
			s += fmt.Sprintf("[%s] %s\n", timeString, msg.Content)
		case LEAVE_CHAT:
			s += fmt.Sprintf("[%s] User %s has left chat %s\n", timeString, msg.SenderID, msg.ChatID)
		case INVITE_TO_CHAT:
			s += fmt.Sprintf("[%s] User %s has been invited by %s\n", timeString, msg.ReceiverID, msg.SenderID)
		case SEND_FILE:
			s += fmt.Sprintf("[%s] User %s has sent a file in chat %s\n", timeString, msg.SenderID, msg.ChatID)
		case SET_USERNAME:
			s += fmt.Sprintf("[%s] %s\n", timeString, msg.Content)
		case LOAD_MESSAGES:
			s += fmt.Sprintf("[%s] %s: %s\n", timeString, msg.SenderID, msg.Content)
		default:
			s += fmt.Sprintf("[%s] Unknown operation received from %s\n", timeString, msg.SenderID)
		}
	}

	s += fmt.Sprintf("\n%s", m.input.View())
	s += "\nPress Enter to send the message, ESC to go back."

	return s
}

// helpView returns the view for the help screen.
func (m Model) helpView() string {
	return `
Available commands in the chat:

  /leave - Leave a chat
  /invite <OnionID> - Invite a user to a chat
  /sendfile <FilePath> - Send a file in a chat (still WIP)
  /setusername <NewUsername> - Set or change the user's username
  /loadmessages - Loads 50 more messages of the chat if they exist

Press ESC to return to the main menu.
`
}

// UsernameView returns the view for the username screen.
func (m Model) UsernameView() string {
	return fmt.Sprintf("Please set your username for %s:\n\n%s\n\nPress Enter to confirm.", m.chatNames[m.CurrentChat], m.usernameInput.View())
}

// TestUserView returns the view for the test user screen.
func (m Model) TestUserView() string {
	return fmt.Sprintf("Test connection to OnionID:\n\n%s\n\nPress Enter to test, ESC to return.", m.testUserInput.View())
}

// CreateChatView returns the view for the create chat screen.
func (m Model) CreateChatView() string {
	s := "Create a new chat:\n\n"
	s += fmt.Sprintf("\nInvitee: %s\n", m.createChatInput.View())
	s += fmt.Sprintf("Chat Name: %s\n", m.chatNameInput.View())
	s += "Invitees:\n"
	for _, invitee := range m.chatInvitees {
		s += fmt.Sprintf("- %s\n", invitee)
	}
	s += "\nPress Enter to add invitee, Ctrl+I to switch input, Ctrl+C to create chat, ESC to return."

	return s
}

// GetSortedChatIDS returns sorted chat IDs.
func GetSortedChatIDs(chatNames map[string]string) []string {
	chatIDs := make([]string, 0, len(chatNames))
	for chatID := range chatNames {
		chatIDs = append(chatIDs, chatID)
	}
	sort.Strings(chatIDs)
	return chatIDs
}

// TestConnection simulates testing a connection to an OnionID.
func TestConnection(onionID string) bool {

	return onionID == "validOnionID"
}

// ClearTempMessage clears the temporary message after a timeout.
func (m Model) ClearTempMessage() tea.Cmd {
	return tea.Tick(time.Until(m.TempMessageExpire), func(time.Time) tea.Msg {
		return TempMsgTimeoutMsg{}
	})
}

// RunFrontend starts the frontend application
func RunFrontend() {
	p := tea.NewProgram(InitialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Oh no, something went wrong: %s\n", err)
		os.Exit(1)
	}
}
