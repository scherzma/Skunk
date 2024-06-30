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

type screen int

const (
	screenIntro screen = iota
	screenChats
	screenHelp
	screenUsername
	screenTestUser
	screenCreateChat
)

type tempMsgTimeoutMsg struct{}

type model struct {
	currentScreen     screen                       // Current screen
	chats             map[string][]FrontendMessage // Map from ChatID to messages
	chatNames         map[string]string            // Map from ChatID to chat names
	invites           []string                     // List of invites
	cursor            int                          // Cursor for selecting chats or invites
	focus             string                       // Can be 'chats' or 'invites'
	currentChat       string                       // Currently selected chat
	inChatDetail      bool                         // Whether we are in chat details
	input             textinput.Model              // User input for messages
	usernameInput     textinput.Model              // User input for setting the username
	usernames         map[string]string            // Map from ChatID to username
	testUserInput     textinput.Model              // User input for testing connection
	createChatInput   textinput.Model              // User input for creating chat (invitees)
	chatNameInput     textinput.Model              // User input for creating chat (chat name)
	chatInvitees      []string                     // List of invitees for the new chat
	tempMessage       string                       // Temporary message
	tempMessageExpire time.Time                    // Expiry time for the temporary message
}

// initialModel initializes the initial state of the application.
func initialModel() model {
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

	return model{
		currentScreen: screenIntro,
		chats: map[string][]FrontendMessage{
			"1": {
				createMessage("1", "Alice", "Bob", "Hello, this is chat 1!", "", "", SEND_MESSAGE),
				createMessage("1", "Bob", "Alice", "Hi Alice!", "", "", SEND_MESSAGE),
			},
			"2": {
				createMessage("2", "Charlie", "Dave", "Welcome to chat 2", "", "", SEND_MESSAGE),
			},
		},
		chatNames:       map[string]string{"1": "Chat 1", "2": "Chat 2"},
		invites:         []string{"Invitation from Alice to Chat 3", "Invitation from Bob to Chat 4"},
		focus:           "chats",
		input:           ti,
		usernameInput:   ui,
		usernames:       make(map[string]string),
		testUserInput:   tu,
		createChatInput: ci,
		chatNameInput:   cn,
		chatInvitees:    []string{},
	}
}

// Init initializes the tea program.
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// createMessage creates a new FrontendMessage.
func createMessage(chatID, senderID, receiverID, content, senderAddress, receiverAddress string, op OperationType) FrontendMessage {
	return FrontendMessage{
		Id:              chatID,
		Timestamp:       time.Now().Unix(),
		Content:         sanitizeInput(content),
		SenderID:        sanitizeInput(senderID),
		ReceiverID:      sanitizeInput(receiverID),
		SenderAddress:   sanitizeInput(senderAddress),
		ReceiverAddress: sanitizeInput(receiverAddress),
		ChatID:          chatID,
		Operation:       op,
	}
}

// isValidUsername checks if the provided username is valid.
func isValidUsername(username string) bool {
	if len(username) > 20 {
		return false
	}
	for _, r := range username {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// sanitizeInput trims whitespace from the input string.
func sanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// validateInput sanitizes the input and checks if it meets the requirements.
func validateInput(input string, maxLength int) (string, error) {
	sanitized := sanitizeInput(input)
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

// handleCommand processes commands entered by the user.
func (m *model) handleCommand(command, argument string) (tea.Model, tea.Cmd) {
	var msg FrontendMessage
	var err error
	switch command {
	case "/leave":
		msg = createMessage(m.currentChat, m.usernames[m.currentChat], "", "", "", "", LEAVE_CHAT)
		m.chats[m.currentChat] = append(m.chats[m.currentChat], msg)
		m.inChatDetail = false
		delete(m.chats, m.currentChat)
		delete(m.chatNames, m.currentChat)
		delete(m.usernames, m.currentChat)
		m.currentChat = ""
	case "/invite":
		argument, err = validateInput(argument, 56)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		msg = createMessage(m.currentChat, m.usernames[m.currentChat], argument, argument, "", "", INVITE_TO_CHAT)
	case "/sendfile":
		argument, err = validateInput(argument, 256)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		msg = createMessage(m.currentChat, m.usernames[m.currentChat], "", argument, "", "", SEND_FILE)
	case "/setusername":
		argument, err = validateInput(argument, 20)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		if isValidUsername(argument) {
			oldUsername := m.usernames[m.currentChat]
			m.usernames[m.currentChat] = argument
			msg = createMessage(m.currentChat, oldUsername, argument, fmt.Sprintf("%s changed their username to %s", oldUsername, argument), "", "", SET_USERNAME)
		} else {
			msg = createMessage(m.currentChat, m.usernames[m.currentChat], "", "Invalid username: "+argument, "", "", SEND_MESSAGE)
		}
	case "/loadmessages":
		msg = createMessage(m.currentChat, m.usernames[m.currentChat], "", "Loading more messages...", "", "", LOAD_MESSAGES)
	default:
		argument, err = validateInput(argument, 256)
		if err != nil {
			return m, tea.Printf("Error: %v", err)
		}
		msg = createMessage(m.currentChat, m.usernames[m.currentChat], "", "Unknown command: "+command, "", "", SEND_MESSAGE)
	}

	m.chats[m.currentChat] = append(m.chats[m.currentChat], msg)
	m.input.SetValue("")
	return m, nil
}

// handleEnter processes the enter key for sending messages or commands.
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
			msg := createMessage(m.currentChat, m.usernames[m.currentChat], "", input, "", "", SEND_MESSAGE)
			m.chats[m.currentChat] = append(m.chats[m.currentChat], msg)
			m.input.SetValue("")
		}
	}
	return m, nil
}

// handleInviteAccept processes the acceptance of a chat invitation.
func (m *model) handleInviteAccept() (tea.Model, tea.Cmd) {
	invitation := m.invites[m.cursor]
	parts := strings.Split(invitation, " to ")
	chatID := parts[len(parts)-1]
	chatName := chatID

	if _, exists := m.chatNames[chatID]; !exists {
		m.chatNames[chatID] = chatName
		m.chats[chatID] = []FrontendMessage{}
	}

	m.invites = append(m.invites[:m.cursor], m.invites[m.cursor+1:]...)
	if len(m.invites) == 0 {
		m.cursor = 0
		m.focus = "chats"
	} else if m.cursor >= len(m.invites) {
		m.cursor = len(m.invites) - 1
	}
	m.tempMessage = fmt.Sprintf("Invitation accepted: %s", invitation)
	m.tempMessageExpire = time.Now().Add(10 * time.Second)
	return m, m.clearTempMessage()
}

// handleInviteDecline processes the decline of a chat invitation.
func (m *model) handleInviteDecline() (tea.Model, tea.Cmd) {
	m.invites = append(m.invites[:m.cursor], m.invites[m.cursor+1:]...)
	if len(m.invites) == 0 {
		m.cursor = 0
		m.focus = "chats"
	} else if m.cursor >= len(m.invites) {
		m.cursor = len(m.invites) - 1
	}
	m.tempMessage = "Invitation declined"
	m.tempMessageExpire = time.Now().Add(10 * time.Second)
	return m, m.clearTempMessage()
}

// handleChatSelection processes the selection of a chat from the list.
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

// Update handles all incoming messages and updates the model accordingly.
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
					joinMsg := createMessage(m.currentChat, username, "", fmt.Sprintf("%s joined the chat", username), "", "", JOIN_CHAT)
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
		case screenTestUser:
			switch msg.Type {
			case tea.KeyEnter:
				onionID := strings.TrimSpace(m.testUserInput.Value())
				m.testUserInput.SetValue("")
				if testConnection(onionID) {
					m.tempMessage = fmt.Sprintf("Connection to %s successful!", onionID)
				} else {
					m.tempMessage = fmt.Sprintf("Failed to connect to %s.", onionID)
				}
				m.tempMessageExpire = time.Now().Add(10 * time.Second)
				return m, m.clearTempMessage()
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
				m.chats[chatID] = []FrontendMessage{}
				for _, invitee := range m.chatInvitees {
					msg := createMessage(chatID, m.usernames[m.currentChat], invitee, invitee, "", "", INVITE_TO_CHAT)
					m.chats[chatID] = append(m.chats[chatID], msg)
				}
				m.createChatInput.SetValue("")
				m.chatNameInput.SetValue("")
				m.chatInvitees = []string{}
				m.currentScreen = screenChats
				m.focus = "chats"
				m.tempMessage = fmt.Sprintf("Chat %s created!", chatName)
				m.tempMessageExpire = time.Now().Add(10 * time.Second)
				return m, m.clearTempMessage()
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
	case tempMsgTimeoutMsg:
		if time.Now().After(m.tempMessageExpire) {
			m.tempMessage = ""
		}
	}

	return m, tea.Batch(cmds...)
}

// handleChatInput processes chat input from the user.
func (m *model) handleChatInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyLeft && len(m.input.Value()) == 0 {
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// toggleFocus switches the focus between chats and invites.
func toggleFocus(focus string) string {
	if focus == "chats" {
		return "invites"
	}
	return "chats"
}

// View renders the current view based on the model's state.
func (m model) View() string {
	tempMsg := ""
	if m.tempMessage != "" && time.Now().Before(m.tempMessageExpire) {
		tempMsg = fmt.Sprintf("\n\n%s\n\n", m.tempMessage)
	}
	switch m.currentScreen {
	case screenIntro:
		return m.introView()
	case screenChats:
		return m.chatsView() + tempMsg
	case screenHelp:
		return m.helpView()
	case screenUsername:
		return m.usernameView()
	case screenTestUser:
		return m.testUserView() + tempMsg
	case screenCreateChat:
		return m.createChatView()
	}
	return ""
}

// introView renders the introductory view.
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

// chatsView renders the main chat list view.
func (m model) chatsView() string {
	if m.inChatDetail {
		return m.chatDetailView()
	}

	s := `
╔══════════════════════════════════════════════════════════════════╗
║                           Skunk                                  ║
║       Secure Peer-to-Peer Communication Platform                 ║
╚══════════════════════════════════════════════════════════════════╝
`

	s += "╔══════════════════════════════════════════════════════════════════╗\n"
	s += "║            Chats                                                 ║\n"
	s += "╠══════════════════════════════════════════════════════════════════╣\n"

	chatIDs := getSortedChatIDs(m.chatNames)

	for i, chatID := range chatIDs {
		cursor := " "
		if i == m.cursor && m.focus == "chats" {
			cursor = ">"
		}
		padding := len(" " + cursor + " " + m.chatNames[chatID])
		spaces := 66 - padding
		s += fmt.Sprintf("║ %s %s%*s║\n", cursor, m.chatNames[chatID], spaces, "")
	}

	s += "╠══════════════════════════════════════════════════════════════════╣\n"
	s += "║           Invites                                                ║\n"
	s += "╠══════════════════════════════════════════════════════════════════╣\n"

	for i, invite := range m.invites {
		cursor := " "
		if i == m.cursor && m.focus == "invites" {
			cursor = ">"
		}
		padding := len(" " + cursor + " " + invite)
		spaces := 66 - padding
		s += fmt.Sprintf("║ %s %s%*s║\n", cursor, invite, spaces, "")
	}

	if m.focus == "invites" && len(m.invites) == 0 {
		s += "║ No invitations available                               ║\n"
	}

	s += "╚══════════════════════════════════════════════════════════════════╝\n"
	s += "\nPress Enter to open the chat, Tab to switch between chats and invites, Backspace to decline an invite, h for help, t to test connection, c to create a chat, and q to quit."

	return s
}

// chatDetailView renders the detailed view of a specific chat.
func (m model) chatDetailView() string {
	s := fmt.Sprintf("Chat: %s\n\n", m.chatNames[m.currentChat])

	for _, msg := range m.chats[m.currentChat] {
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

// helpView renders the help screen.
func (m model) helpView() string {
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

// usernameView renders the view for setting the username.
func (m model) usernameView() string {
	return fmt.Sprintf("Please set your username for %s:\n\n%s\n\nPress Enter to confirm.", m.chatNames[m.currentChat], m.usernameInput.View())
}

// testUserView renders the view for testing a connection to an OnionID.
func (m model) testUserView() string {
	return fmt.Sprintf("Test connection to OnionID:\n\n%s\n\nPress Enter to test, ESC to return.", m.testUserInput.View())
}

// createChatView renders the view for creating a new chat.
func (m model) createChatView() string {
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

// getSortedChatIDs returns a sorted slice of chat IDs.
func getSortedChatIDs(chatNames map[string]string) []string {
	chatIDs := make([]string, 0, len(chatNames))
	for chatID := range chatNames {
		chatIDs = append(chatIDs, chatID)
	}
	sort.Strings(chatIDs)
	return chatIDs
}

// testConnection simulates a connection test to the OnionID.
func testConnection(onionID string) bool {
	// Simulate a connection test to the OnionID
	// Replace with actual connection logic as needed
	return onionID == "validOnionID"
}

// clearTempMessage clears the temporary message after a specified duration.
func (m model) clearTempMessage() tea.Cmd {
	return tea.Tick(time.Until(m.tempMessageExpire), func(time.Time) tea.Msg {
		return tempMsgTimeoutMsg{}
	})
}

// main initializes and starts the tea program.
func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Oh no, something went wrong: %s\n", err)
		os.Exit(1)
	}
}
