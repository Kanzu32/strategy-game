package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type Database struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewDatabase(uri, dbName, collectionName string) *Database {
	log.Printf("Подключение к MongoDB по URI: %s", uri)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Не удалось проверить подключение к MongoDB: %v", err)
	}

	log.Printf("Успешное подключение к БД %s, коллекция %s", dbName, collectionName)
	collection := client.Database(dbName).Collection(collectionName)
	return &Database{client: client, collection: collection}
}

func (db *Database) Authenticate(username, password string) bool {
	log.Printf("Аутентификация пользователя: %s", username)
	filter := bson.M{"username": username}
	var result Account
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Printf("Ошибка поиска пользователя: %v", err)
		return false
	}

	if result.Password != password {
		log.Printf("Неверный пароль для пользователя %s", username)
		return false
	}

	log.Printf("Успешная аутентификация пользователя %s", username)
	return true
}

func (db *Database) Register(account *Account) bool {
	log.Printf("Регистрация нового пользователя: %s", account.Username)
	filter := bson.M{"username": account.Username}
	var existing Account
	err := db.collection.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		log.Printf("Пользователь %s уже существует", account.Username)
		return false
	}

	_, err = db.collection.InsertOne(context.TODO(), account)
	if err != nil {
		log.Printf("Ошибка при регистрации пользователя: %v", err)
		return false
	}

	log.Printf("Пользователь %s успешно зарегистрирован", account.Username)
	return true
}
