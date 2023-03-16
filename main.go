package main

import (
	"martini/controllers"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
)

func main() {
	m := martini.Classic()
	//Bisa ada 2 jenis middleware, yang global (artinya dipakai di semua router)
	//dan yang local(bisa di satu grup ataupun di satu router)
	//contoh middleware : CORS
	//cors saat dipakai global
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*.foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	allowCORSHandler := cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*.foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
	//cors saat dipakai di group
	m.Group("/user", func(r martini.Router) {
		r.Get("/", controllers.Authenticate(controllers.GetAllUsers, 1))
		r.Post("/", controllers.Authenticate(controllers.InserNewUser, 2))
		r.Put("/:id", controllers.UpdateUser)
		r.Delete("/:id", controllers.DeleteUser)
	}, allowCORSHandler)
	//cors saat dipakai individu
	m.Put("/corsUser", allowCORSHandler, controllers.GetAllUsers)

	m.Get("/login", controllers.Login)
	m.Get("/logout", controllers.Logout)
	//untuk running di port 8080
	m.RunOnAddr(":8080")
}
