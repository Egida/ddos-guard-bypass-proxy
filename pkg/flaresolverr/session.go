package flaresolverr

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type Session struct {
	c    Client
	lock *Lock
}

func NewSession(url string) *Session {
	s := &Session{
		c: Client{
			url: url,
		},
	}

	s.lock = NewLock(10*time.Second, s.destroySession)

	return s
}

func (s *Session) Get(ctx context.Context, url string) (*http.Response, error) {
	ok := s.lock.TryLock()
	if !ok {
		return nil, errors.New("unable to lock")
	}
	defer s.lock.UnLock()

	sessionID, err := s.sessionID(ctx, s.c.SessionID())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s.c.SetSessionID(sessionID)

	res, err := s.c.Get(ctx, url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}

func (s *Session) sessionID(
	ctx context.Context,
	oldSessionID string,
) (string, error) {
	sessions, err := s.c.ListSessions(ctx)
	if err != nil {
		return "", err
	}

	for _, id := range sessions {
		if id == oldSessionID {
			continue
		}

		err := s.c.DestroySession(ctx, id)
		if err != nil {
			return "", err
		}
	}

	if oldSessionID != "" && contains(sessions, oldSessionID) {
		return oldSessionID, nil
	}

	return s.c.CreateSession(ctx)
}

func (s *Session) destroySession() {
	log.Println("auto session destroy")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.c.DestroySession(ctx, s.c.SessionID())
	if err != nil {
		log.Println(err)
	}

	s.c.SetSessionID("")
}

func contains[K comparable](s []K, i K) bool {
	for _, x := range s {
		if x == i {
			return true
		}
	}

	return false
}
