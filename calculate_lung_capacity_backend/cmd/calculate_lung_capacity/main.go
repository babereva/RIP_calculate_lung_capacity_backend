package main

import (
	"log"

	"calculate_lung_capacity/internal/api"
)

func main() {
	log.Println("Lung capacity service start!")
	api.StartServer()
	log.Println("Lung capacity service terminated!")
}