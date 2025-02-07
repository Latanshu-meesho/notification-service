package utils

import (
	"log"
	"os"
)

func InitializeLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Logger initialized.")
}
