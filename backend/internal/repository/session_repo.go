package repository

import (
	"errors"
	"time"

	"pdf-parser/internal/model"
	"pdf-parser/pkg/database"
)

type SessionRepository struct{}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (r *SessionRepository) CreateSession(session *model.Session) error {
	return database.GetDB().Create(session).Error
}

func (r *SessionRepository) GetSession(sessionID string) (*model.Session, error) {
	var session model.Session
	err := database.GetDB().Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetSessionByAccessToken(accessToken string) (*model.Session, error) {
	var session model.Session
	err := database.GetDB().Where("access_token = ?", accessToken).First(&session).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) UpdateSession(session *model.Session) error {
	return database.GetDB().Save(session).Error
}

func (r *SessionRepository) DeleteSession(sessionID string) error {
	return database.GetDB().Where("session_id = ?", sessionID).Delete(&model.Session{}).Error
}

func (r *SessionRepository) DeleteExpiredSessions() error {
	return database.GetDB().Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}

func (r *SessionRepository) SaveCodeExchange(state, codeVerifier, redirectURI string) error {
	codeExchange := &model.CodeExchange{
		State:        state,
		CodeVerifier: codeVerifier,
		RedirectURI:  redirectURI,
		Used:         false,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}
	return database.GetDB().Create(codeExchange).Error
}

func (r *SessionRepository) GetCodeExchange(state string) (*model.CodeExchange, error) {
	var codeExchange model.CodeExchange
	err := database.GetDB().Where("state = ? AND used = ? AND expires_at > ?", state, false, time.Now()).First(&codeExchange).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, errors.New("code exchange not found or expired")
		}
		return nil, err
	}
	return &codeExchange, nil
}

func (r *SessionRepository) MarkCodeExchangeUsed(state string) error {
	return database.GetDB().Model(&model.CodeExchange{}).Where("state = ?", state).Update("used", true).Error
}
