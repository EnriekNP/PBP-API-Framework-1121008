package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response ErrorResponse
	response.Status = 401
	response.Message = "Unauthorized Access"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func sendErrorResponse(w http.ResponseWriter, message string) {
	var response ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// function to Login
func Login(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()
	//Read from header
	platform := r.Header.Get("platform")

	//Read From Query Param
	email := r.URL.Query()["Email"]
	password := r.URL.Query()["Password"]

	query := "SELECT ID,Name,Age,Address,UserType FROM USERS WHERE Email ='" + email[0] + "' && Password='" + password[0] + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var user User
	login := false
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.UserType); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			generateToken(w, user.ID, user.Name, user.UserType)
			login = true
		}
	}
	var response UserResponse
	if login {
		response.Status = 200
		response.Message = "Success login from " + platform
		response.Data = user

	} else {
		response.Status = 400
		response.Message = "Login Failed"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// function to logout
func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	var response UserResponse
	response.Status = 200
	response.Message = "Success"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fungsi untuk menampilkan semua user
// atau untuk menampilkan semua user yang difilter berdasarkan nama atau umur atau alamat
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	fmt.Println("Masuk GetAllUsers")
	query := "SELECT * FROM USERS"

	//membaca filter dari query param
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]
	address := r.URL.Query()["address"]
	//mengecek apakah ada query param
	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name='" + name[0] + "'"
	}
	if age != nil {
		if name != nil {
			query += " AND"
		} else {
			query += " WHERE"
		}
		fmt.Println(age[0])
		query += " age ='" + age[0] + "'"
	}
	if address != nil {
		if name != nil || age != nil {
			query += " AND"
		} else {
			query += " WHERE"
		}
		fmt.Println(address[0])
		query += " address LIKE'%" + address[0] + "%'"
	}

	//eksekusi query
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	//memasukan hasil query ke dalam variable tipe user
	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			users = append(users, user)
		}
	}
	//kirim response
	w.Header().Set("Content-Type", "application/json")
	var response UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

// Insert New User
func InserNewUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	userType, _ := strconv.Atoi(r.Form.Get("usertype"))
	//execute the query to insert into database
	res, errQuery := db.Exec("INSERT INTO users(name,age,address,email,password,UserType)values(?,?,?,?,?,?)",
		name,
		age,
		address,
		email,
		password,
		userType,
	)
	//variable to get Id users
	id, _ := res.LastInsertId()
	//add all variable into user struct
	var user User
	user.ID = int(id)
	user.Name = name
	user.Age = age
	user.Address = address
	user.Email = email
	user.Password = password
	user.UserType = userType
	//return response , if there are no error then return response with message sucess and data user,
	//else return response with failed message
	var response UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = user
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Insert User Failed"
	}
	//to make response into json type and return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// function to update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	id, _ := strconv.Atoi(r.Form.Get("id"))
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	_, errQuery := db.Exec("UPDATE users SET name =?,age = ?,address = ? WHERE ID=?",
		name,
		age,
		address,
		id,
	)

	var response UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data.ID = id
		response.Data.Name = name
		response.Data.Age = age
		response.Data.Address = address
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Update User Failed"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// delete user
func DeleteUser(params martini.Params, w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	userId  := params["id"]

	_, errQuery := db.Exec("DELETE FROM users WHERE id=?",
	userid)

	var response UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Delete Failed"
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}