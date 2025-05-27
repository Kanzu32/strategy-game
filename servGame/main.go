package main

import (
	"SERV/database"
	"SERV/server"
	"SERV/terminal"
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port      string
	MongoURI  string
	DbName    string
	CollUsers string
	CollStats string
	CollMaps  string
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "PORT":
			config.Port = value
		case "MONGO_URI":
			config.MongoURI = value
		case "DB_NAME":
			config.DbName = value
		case "COLL_USERS":
			config.CollUsers = value
		case "COLL_STATS":
			config.CollStats = value
		case "COLL_MAPS":
			config.CollMaps = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	terminal.Start()

	terminal.Log("Initializing game server...")

	config, err := LoadConfig("config.txt")
	if err != nil {
		terminal.LogFatal("Failed to load config:", err)
	}

	terminal.Log("Connecting to database...")
	database := database.NewDatabase(
		config.MongoURI,
		config.DbName,
		config.CollUsers,
		config.CollStats,
		config.CollMaps,
	)

	terminal.Log("Starting server...")
	server := server.NewServer(database)
	server.Start(config.Port)
}
