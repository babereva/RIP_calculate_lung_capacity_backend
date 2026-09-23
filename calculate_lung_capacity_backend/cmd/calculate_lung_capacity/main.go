package main

import (
	"log"

	"calculate_lung_capacity/internal/api"
)

func main() {
	log.Println("Patient category service start!")
	api.StartServer()
	log.Println("Patient category service terminated!")
}