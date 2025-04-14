package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Config struct {
	Port     string
	MongoURI string
	DbName   string
	CollName string
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
		case "COLL_NAME":
			config.CollName = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	log.Println("Initializing game server...")

	config, err := LoadConfig("config.txt")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Connecting to database...")
	database := NewDatabase(config.MongoURI, config.DbName, config.CollName)

	log.Println("Starting server...")
	server := NewServer(database)
	server.Start(config.Port)
}
