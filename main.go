package main

import "github.com/nikola43/fibergormapitemplate/app"

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:3001
// @BasePath /api/v1
func main() {
	a := new(app.App)
	a.Initialize(":3001")
}
