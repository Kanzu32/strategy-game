package terminal

import (
	"log"
	"os"
	"time"
)

var fileLoger log.Logger
var file os.File

func Start() {
	filepath := "logs/" + time.Now().Format("2006-01-02 15-04-05 MST") + ".txt"
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalln("Ошибка при создании файла логирования", err)
	}
	fileLoger = *log.New(file, "", log.Ldate|log.Ltime)
	Log("Терминал управления сервером запущен.")
}

func Shutdown() {
	Log("Выключение сервера...")
	file.Close()
	os.Exit(0)
}

func Log(v ...any) {
	fileLoger.Println(v...)
	log.Println(v...)
}

func LogFatal(v ...any) {
	Log(v...)
	Shutdown()
}
