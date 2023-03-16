package main

import (
	"martini/controllers"

	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Group("/user", func(r martini.Router) {
		r.Get("/", controllers.GetAllUsers)
		r.Post("/", controllers.InserNewUser)
		r.Put("/:id", controllers.UpdateUser)
		r.Delete("/:id", controllers.DeleteUser)
	})

	m.Get("/login", controllers.Login)
	m.Get("/logout", controllers.Logout)
	m.RunOnAddr(":8080")
}
