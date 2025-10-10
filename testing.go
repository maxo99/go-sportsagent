package main

import "github.com/joho/godotenv"

func init() {
	// Load .env file automatically for all tests in the main package
	godotenv.Load(".env")
}
