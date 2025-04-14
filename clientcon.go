package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ClientConnection управляет подключением между сервером и клиентом через HTTP.
type ClientConnection struct {
	w http.ResponseWriter
	r *http.Request
}

// NewClientConnection создает новое клиентское соединение.
func NewClientConnection(w http.ResponseWriter, r *http.Request) *ClientConnection {
	return &ClientConnection{w: w, r: r}
}

// SendData отправляет данные клиенту в формате JSON.
func (c *ClientConnection) SendData(data string) error {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(http.StatusOK)
	_, err := c.w.Write([]byte(data))
	return err
}

// ReceiveData получает данные от клиента.
func (c *ClientConnection) ReceiveData() (map[string]string, error) {
	body, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return nil, err
	}

	var request map[string]string
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, err
	}

	return request, nil
}
