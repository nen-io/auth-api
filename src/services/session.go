package services

type Session struct {
	State         map[string]bool
	BlockedTokens map[string]bool
}

var SessionManager *Session

func InitSession() {
	SessionManager = &Session{State: make(map[string]bool)}
}

func (s *Session) Add(id string) {
	s.State[id] = true
}
