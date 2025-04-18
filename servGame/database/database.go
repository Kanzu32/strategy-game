package database

import (
	"context"
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

func (db *Database) Authenticate(email, password string) bool {
	log.Printf("Аутентификация пользователя: %s", email)
	filter := bson.M{"email": email}
	var result Account
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Printf("Ошибка поиска пользователя: %v", err)
		return false
	}

	if result.Password != password {
		log.Printf("Неверный пароль для пользователя %s", email)
		return false
	}

	log.Printf("Успешная аутентификация пользователя %s", email)
	return true
}

func (db *Database) Register(account *Account) bool {
	log.Printf("Регистрация нового пользователя: %s", account.Email)
	filter := bson.M{"email": account.Email}
	var existing Account
	err := db.collection.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		log.Printf("Пользователь %s уже существует", account.Email)
		return false
	}

	_, err = db.collection.InsertOne(context.TODO(), account)
	if err != nil {
		log.Printf("Ошибка при регистрации пользователя: %v", err)
		return false
	}

	log.Printf("Пользователь %s успешно зарегистрирован", account.Email)
	return true
}
