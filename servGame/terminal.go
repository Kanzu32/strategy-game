package main

import "log"

type Terminal struct{}

func NewTerminal() *Terminal {
	return &Terminal{}
}

func (t *Terminal) Start(server *Server) {
	log.Println("Терминал управления сервером запущен.")
	// Логика взаимодействия с администратором
}

func (t *Terminal) Shutdown(server *Server) {
	log.Println("Выключение сервера...")
	// Логика завершения работы сервера
}
