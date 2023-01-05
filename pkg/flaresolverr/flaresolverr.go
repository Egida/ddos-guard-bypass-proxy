package flaresolverr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var ErrFlareSolverr = errors.New("FlareSolverr error")

type cmd = int

const (
	get cmd = iota
	sessionsCreate
	sessionsDestroy
	sessionsList
)

func formatCMD(c cmd) string {
	switch c {
	case get:
		return "request.get"
	case sessionsCreate:
		return "sessions.create"
	case sessionsDestroy:
		return "sessions.destroy"
	case sessionsList:
		return "sessions.list"
	default:
		log.Panicln(errors.New("unknown cmd"))
		return ""
	}
}

type query struct {
	CMD            string `json:"cmd"`
	URL            string `json:"url,omitempty"`
	MaxTimeoutInMS int    `json:"maxTimeout,omitempty"`
	SessionID      string `json:"session,omitempty"`
}

type Client struct {
	Url       string
	SessionID string
}

func (c *Client) do(ctx context.Context, query query) (*http.Response, error) {
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	return client.Do(req)
}

func (c *Client) ListSessions(ctx context.Context) ([]string, error) {
	query := query{
		CMD: formatCMD(sessionsList),
	}

	res, err := c.do(ctx, query)
	if err != nil {
		return nil, err
	}

	return parseListSessionsResponse(res)
}

func (c *Client) CreateSession(ctx context.Context) (string, error) {
	query := query{
		CMD: formatCMD(sessionsCreate),
	}

	res, err := c.do(ctx, query)
	if err != nil {
		return "", err
	}

	return parseCreateSessionResponse(res)
}

func (c *Client) DestroySession(
	ctx context.Context,
	sessionID string,
) error {
	query := query{
		CMD:       formatCMD(sessionsDestroy),
		SessionID: sessionID,
	}

	res, err := c.do(ctx, query)
	if err != nil {
		return err
	}

	return parseDestroySessionResponse(res)
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	query := query{
		CMD:       formatCMD(get),
		URL:       url,
		SessionID: c.SessionID,
	}

	res, err := c.do(ctx, query)
	if err != nil {
		return nil, err
	}

	return parseGetResponse(res)
}

type listSessionsResponse struct {
	Status         string   `json:"status"`
	Message        string   `json:"message"`
	StartTimestamp int64    `json:"startTimestamp"`
	EndTimestamp   int64    `json:"endTimestamp"`
	Version        string   `json:"version"`
	Sessions       []string `json:"sessions"`
}

func parseListSessionsResponse(res *http.Response) ([]string, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var payload listSessionsResponse
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Status != "ok" {
		err := fmt.Errorf(
			"%w: %s",
			ErrFlareSolverr,
			payload.Message,
		)
		return nil, err
	}

	sessions := payload.Sessions

	return sessions, nil
}

type createSessionResponse struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	Version        string `json:"version"`
	Session        string `json:"session"`
}

func parseCreateSessionResponse(res *http.Response) (string, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var payload createSessionResponse
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return "", err
	}

	if payload.Status != "ok" {
		err := fmt.Errorf(
			"%w: %s",
			ErrFlareSolverr,
			payload.Message,
		)
		return "", err
	}

	session := payload.Session

	return session, nil
}

type destroySessionResponse struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	Version        string `json:"version"`
}

func parseDestroySessionResponse(res *http.Response) error {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var payload destroySessionResponse
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return err
	}

	if payload.Status != "ok" {
		err := fmt.Errorf(
			"%w: %s",
			ErrFlareSolverr,
			payload.Message,
		)
		return err
	}

	return nil
}

type getResponse struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	Version        string `json:"version"`
	Solution       struct {
		URL       string `json:"url"`
		Status    int    `json:"status"`
		Headers   map[string]string
		Response  string        `json:"response"`
		Cookies   []interface{} `json:"cookies"`
		UserAgent string        `json:"userAgent"`
	} `json:"solution"`
}

func parseGetResponse(res *http.Response) (*http.Response, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var payload getResponse
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Status != "ok" {
		err := fmt.Errorf(
			"%w: %s",
			ErrFlareSolverr,
			payload.Message,
		)
		return nil, err
	}

	header := map[string][]string{}
	for k, v := range payload.Solution.Headers {
		header[k] = []string{v}
	}

	body := io.NopCloser(strings.NewReader(payload.Solution.Response))

	r := &http.Response{
		StatusCode: payload.Solution.Status,
		Header:     header,
		Body:       body,
	}

	return r, nil
}
