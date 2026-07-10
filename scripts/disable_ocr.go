package main

import (
	"fmt"
	"log"

	"harmony_api/config"
)

func main() {
	config.LoadConfig()
	defer config.CloseDB()

	fmt.Println("Disabling analyze_incoming_images for all channel integrations to stabilize production...")
	res := config.DB.Exec("UPDATE channel_integrations SET analyze_incoming_images = false")
	if res.Error != nil {
		log.Fatalf("Error updating database: %v\n", res.Error)
	}
	fmt.Printf("SUCCESS! Rows affected: %d\n", res.RowsAffected)
}
