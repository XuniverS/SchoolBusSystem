package main

import "awesomeProject/back"

func main() {

	router := back.SetupRoutes()
	router.Run(":8000")

}
