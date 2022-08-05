package main

import (
	"log"
	config "team3-task/config"
	"team3-task/internal/app"
)

// @title Team3.Task.service API
// @version 1.0
// @description This is a team3 task.service API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://gitlab.com/g6834/team3/task
// @contact.email alex@mail.ru

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host :3000
// @BasePath /
// @query.collection.format multi

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
