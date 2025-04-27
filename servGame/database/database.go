package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
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

func (db *Database) Authenticate(email, password string) error {
	// log.Printf("Аутентификация пользователя: %s", email)
	filter := bson.M{"email": email}
	var result Account
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка поиска пользователя в базе данных: %v", err))
	}

	if result.Password != password {
		// log.Printf("Неверный пароль для пользователя %s", email)
		return errors.New(fmt.Sprintf("Неверный пароль для пользователя %s", email))
	}

	// log.Printf("Успешная аутентификация пользователя %s", email)
	return nil
}

func (db *Database) Register(account *Account) error {
	// log.Printf("Регистрация нового пользователя: %s", account.Email)
	filter := bson.M{"email": account.Email}
	var existing Account
	err := db.collection.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		// log.Printf("Пользователь %s уже существует в базе данных", account.Email)

		return errors.New(fmt.Sprintf("Пользователь %s уже существует в базе данных", account.Email))
	}

	_, err = db.collection.InsertOne(context.TODO(), account)
	if err != nil {
		// log.Printf("Ошибка при добавлении пользователя в базу данных: %v", err)
		return errors.New(fmt.Sprintf("Ошибка при добавлении пользователя в базу данных: %v", err))
	}

	// log.Printf("Пользователь %s успешно добавлен в базу данных", account.Email)
	return nil
}
