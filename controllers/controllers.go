package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
)

// fungsi untuk mengirimkan response akses yang tidak diizinkan
func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response ErrorResponse
	response.Status = 401
	response.Message = "Unauthorized Access"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fungsi untuk mengirimkan response error
func sendErrorResponse(w http.ResponseWriter, message string) {
	var response ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fungsi untuk login
func Login(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()
	//membaca platform dari header
	platform := r.Header.Get("platform")

	//membaca email dan password dari query param
	email := r.URL.Query()["Email"]
	password := r.URL.Query()["Password"]
	//eksekusi query
	query := "SELECT ID,Name,Age,Address,email,password,UserType FROM USERS WHERE Email ='" + email[0] + "' && Password='" + password[0] + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	//masukan ke dalam variable user jika ketemu user dengan email dan password yang benar
	var user User
	login := false
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			//jika ketemu, maka akan mengenerate token dan login menjadi true
			GenerateToken(w, user.ID, user.Name, user.UserType)
			login = true
		}
	}
	//kirimkan response dengan mengecek apakah login berhasil
	//jika berhasil, maka akan mengirimkan data user , status dan pesan berhasil
	//jika gagal, maka akan mengirimkan status dan pesan gagal
	var response UserResponse
	if login {
		response.Status = 200
		response.Message = "Success login from " + platform
		response.Data = user

	} else {
		response.Status = 400
		response.Message = "Login Failed"
	}
	//mengembalikan response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fungsi utnuk logout
func Logout(w http.ResponseWriter, r *http.Request) {
	//menghilangkan token
	resetUserToken(w)
	// mengembalikan response sukses
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

	//membaca dari Request Body
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
	//eksekusi kueri
	res, errQuery := db.Exec("INSERT INTO users(name,age,address,email,password,UserType)values(?,?,?,?,?,?)",
		name,
		age,
		address,
		email,
		password,
		userType,
	)
	//variabel untuk mendapatkan id user
	id, _ := res.LastInsertId()
	//menambahkan semua variabel kedalam struct user
	var user User
	user.ID = int(id)
	user.Name = name
	user.Age = age
	user.Address = address
	user.Email = email
	user.Password = password
	user.UserType = userType
	//megnembalikan respon dengan mengecek apakah ada eror
	//jika tidak ada error akan mengembalikan data user , status dan pesan sukses
	//jika ada error maka akan mengembalikan status dan pesan error
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
	//mereturn respon
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fungsi untuk mengupdate user
func UpdateUser(param martini.Params, w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	//membaca dari request body
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	//membaca id dari param
	id := param["id"]
	//mengkueri terlebih dahulu dari database user dengan id ini
	query := "SELECT name,age,address,email,password,UserType FROM USERS WHERE Id ='" + id + "'"
	//eksekusi kueri select user
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	//memasukan user tersebut ke variable user
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		}
	}
	//membaca request body
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	userType, _ := strconv.Atoi(r.Form.Get("userType"))
	//mengecek apakah request body yang didapat itu kosong atau tidak
	//jika kosong maka akan diisi oleh data lama yang telah diambil dari database
	if name == "" {
		name = user.Name
	}
	if age == 0 {
		age = user.Age
	}
	if address == "" {
		address = user.Address
	}
	if email == "" {
		email = user.Email
	}
	if password == "" {
		password = user.Password
	}

	if userType == 0 {
		userType = user.UserType
	}
	//mengeksekusi query update user
	_, errQuery := db.Exec("UPDATE users SET name =?,age = ?,address = ?,email= ?,password= ?,UserType= ? WHERE ID=?",
		name,
		age,
		address,
		email,
		password,
		userType,
		id,
	)
	//menambahkan user akhir yang sudah pasti ke dalam variabel user final
	var userFinal User
	userFinal.ID, _ = strconv.Atoi(id)
	userFinal.Name = name
	userFinal.Age = age
	userFinal.Address = address
	userFinal.Email = email
	userFinal.Password = password
	userFinal.UserType = userType
	//mengecek apakah ada error di query
	//jika tidak ada maka akan menampilkan data user , pesan dan status sukses
	var response UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = userFinal
		GenerateToken(w, user.ID, user.Name, user.UserType)
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Update User Failed"
	}
	//mereturn response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// delete user
func DeleteUser(params martini.Params, w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()
	//membaca dari params
	userid := params["id"]
	//eksekusi delete
	_, errQuery := db.Exec("DELETE FROM users WHERE id=?",
		userid)
	//menghasilkan response
	//jika tidak ada eror, maka akan ditampilkan status dan pesan sukses
	// jika ada eror, maka akan diprint erornya dan ditampilkan status dan pesan gagal
	var response UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Delete Failed"
	}
	//untuk mereturn response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
