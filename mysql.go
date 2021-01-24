//
//  go-unit-test-sql
//
//  Copyright © 2020. All rights reserved.
//

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)
type Employee struct {
	ID int `json:"id"`
	Name string `json:"name"`
	RoleID int `json:"roleId"`
	Gender string `json:"gender"`
}
// repository represent the repository model
type sqlDb struct {
	db *sql.DB
}

//For Connection with Database
func NewConnection() (*sqlDb, error) {
	db, err := sql.Open("mysql", "sunny:Sunny@9570@(127.0.0.1)/company")
	if err != nil {
		return nil, err
	}

	return &sqlDb{db}, nil
}

// Close the connection
func (r *sqlDb) Close() {
	r.db.Close()
}

var emp Employee

func  (DB *sqlDb) insertIntoEmployee(empData Employee) int  {

	db := DB.db
	ename,roleId,gender := empData.Name,empData.RoleID,empData.Gender

	//sql := "INSERT INTO Employee(id,name, roleId,gender) VALUES ( %v,%v, %v,%v)"
	sq := "INSERT INTO Employee(name, roleId,gender) VALUES (?,?,?)"

	res, err := db.Exec(sq,ename,roleId,gender)

	if err != nil {
		panic(err.Error())
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)
	return int(lastId)
}


func (DB *sqlDb) selectData() ([]Employee,error){

	db := DB.db

	dis,err := db.Query("SELECT id, name, roleId, gender FROM Employee")

	defer dis.Close()
	var empData [] Employee
	if err != nil {
		log.Fatal(err)
		return empData,err
	}

	for dis.Next() {

		var em Employee
		err := dis.Scan(&em.ID,&em.Name,&em.RoleID,&em.Gender)

		if err != nil {
			log.Fatal(err)
		}
		empData = append(empData, em)
	}
	return empData,nil
}

func (DB *sqlDb) selectByroleId(roleId int) [] Employee{
	db := DB.db
	dis,err := db.Query("SELECT id, name, roleId, gender FROM Employee WHERE roleId = ?",roleId)

	defer dis.Close()
	var empData [] Employee
	if err != nil {
		log.Fatal(err)
	}

	for dis.Next() {

		var em Employee
		err := dis.Scan(&em.ID,&em.Name,&em.RoleID,&em.Gender)

		if err != nil {
			log.Fatal(err)
		}
		empData = append(empData, em)
	}
	return empData
}


func (dbHandler *sqlDb) selectInsert(w http.ResponseWriter,r *http.Request)  {

	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		keys, ok := r.URL.Query()["roleId"]
		//log.Println(r.URL.Query())
		if !ok {
			empData,err := dbHandler.selectData()
			if err!=nil{
				http.Error(w,"DataBase Connection Problem",http.StatusInternalServerError)
				return
			}
			if len(empData) == 0 {
				http.Error(w,"This id employee is not available",http.StatusBadRequest)
				return
			}
			res, er := json.Marshal(empData)
			if er == nil {
				w.Write(res)
			}
			return
		}
		key, er := strconv.Atoi(keys[0])
		if er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}
		empData := dbHandler.selectByroleId(key)
		if len(empData) == 0 {
			http.Error(w,fmt.Sprintf("there is no employee with roleId = %v",key),http.StatusBadRequest)
			return
		}
		res, er := json.Marshal(empData)
		if er == nil {
			w.Write(res)
		}
	}
	if r.Method == "POST" {
		body := r.Body

		err := json.NewDecoder(body).Decode(&emp)

		if err != nil {
			log.Print(err)
		}
		lastId := dbHandler.insertIntoEmployee(emp)
		emp,err := dbHandler.selectIndexData(lastId)
		res,er :=json.Marshal(emp)
		if er == nil {
			w.Write(res)
		}
	}
}

func (DB *sqlDb) selectIndexData(id int) (*Employee,error){
	db := DB.db
	dis,err := db.Query("SELECT id, name, roleId, gender FROM Employee WHERE id = ?",id)
   // fmt.Println(err)
	defer dis.Close()

	if err != nil {
		log.Fatal(err)
	}
	var em Employee
	if dis.Next() {
		err := dis.Scan(&em.ID,&em.Name,&em.RoleID,&em.Gender)

		if err != nil {
			log.Fatal(err)
		}
		return &em,nil
	} else {
		return nil,errors.New("Not Found")
	}
}

func (dbHandler *sqlDb) employeeId(w http.ResponseWriter,r *http.Request) {

	if r.Method == "GET" || r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		v := params["id"]

		id,er:= strconv.Atoi(v)

		//fmt.Println(id)
		if er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}
		empData,err := dbHandler.selectIndexData(id)
		//fmt.Println(flag)
		if err!=nil {
			http.Error(w,fmt.Sprintf("Employee with id %v is not available!",id),http.StatusBadRequest)
			return
		}
		res, er := json.Marshal(*empData)
		if er == nil {
			w.Write(res)
		}

	}

}

func (DB *sqlDb) updatEmployee(empData Employee,id int) int{
	db := DB.db


	ename,roleId,gender := empData.Name,empData.RoleID,empData.Gender

	//sq := fmt.Sprintf("UPDATE Employee SET name = '%s',roleId = %v,gender = '%s' where id = %v",ename,roleId,gender,id)
	//log.Println(sq)
	res, err := db.Exec("UPDATE Employee SET name = ?,roleId = ?,gender = ? where id = ?",ename,roleId,gender,id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	fmt.Printf("Number of Effected Row: %d\n", rowsAffected)
	return int(rowsAffected)
}

func (dbHandler *sqlDb) update(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	v := params["id"]
	id,er:= strconv.Atoi(v)

	//log.Println(id)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "PUT" {
		body := r.Body
		err := json.NewDecoder(body).Decode(&emp)

		if err != nil {
			log.Print(err)
		}
		rowsAffected := dbHandler.updatEmployee(emp,id)
		if rowsAffected==0 {
			http.Error(w,"No records Affected",http.StatusBadRequest)
			return
		}

		emp,err := dbHandler.selectIndexData(id)
		res,er :=json.Marshal(&emp)
		if er == nil {
			w.Write(res)
		}
	}

}

func main(){
	dbHandler,err := NewConnection()
	defer dbHandler.Close()
	if err != nil{
		panic(err)
		return
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/Employee",dbHandler.selectInsert)
	r.HandleFunc("/Employee/{id}",dbHandler.employeeId)
	r.HandleFunc("/Employee/Update/{id}",dbHandler.update).Methods("PUT")
	http.ListenAndServe(":8080",r)
}



