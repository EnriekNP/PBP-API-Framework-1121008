package main

import (
	"martini/controllers"

	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Group("/user", func(r martini.Router) {
		r.Get("/", controllers.Authenticate(controllers.GetAllUsers, 1))
		r.Post("/", controllers.Authenticate(controllers.InserNewUser, 2))
		r.Put("/:id", controllers.UpdateUser)
		r.Delete("/:id", controllers.DeleteUser)
	})
	//Bisa ada 2 jenis middleware, yang global dan yang local
	//contoh CORS

	m.Get("/login", controllers.Login)
	m.Get("/logout", controllers.Logout)
	m.RunOnAddr(":8080")
}
