package main

import (
	"local/cmd"
	_ "local/docs" // swagger docs
)

// @title           Simple Chat API
// @version         1.0
// @description     This is a simple chat application API server.

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:80
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cmd.Execute()
}
