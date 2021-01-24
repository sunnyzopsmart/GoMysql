package main

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var u = []Employee{{
	ID:     1,
	Name:   "Sunny",
	RoleID: 101,
	Gender: "M",
},
	{
		ID:     2,
		Name:   "Rahul",
		RoleID: 101,
		Gender: "M",
	},
}
func TestFindByID(t *testing.T) {
	db, mock,err := sqlmock.New()
	dbHandler := &sqlDb{db}
	defer func() {
		dbHandler.Close()
	}()

	query := "SELECT id, name, roleId, gender FROM Employee WHERE id = ?"

	rows := sqlmock.NewRows([]string{"id", "name", "roleId", "gender"}).
		AddRow(u[0].ID, u[0].Name, u[0].RoleID, u[0].Gender)

	mock.ExpectQuery(query).WithArgs(u[0].ID).WillReturnRows(rows)
     //mock.ExpectQuery(query).WithArgs(1).WillReturnError(err)
	user,err := dbHandler.selectIndexData(u[0].ID)

	emp := *user
	assert.Equal(t, u[0],emp)
	assert.NoError(t, err)
}

func TestSelect(t *testing.T) {
	db, mock,err := sqlmock.New()
	dbHandler := &sqlDb{db}
	defer func() {
		dbHandler.Close()
	}()

	query := "SELECT id, name, roleId, gender FROM Employee"

	rows := sqlmock.NewRows([]string{"id", "name", "roleId", "gender"}).
		AddRow(u[0].ID, u[0].Name, u[0].RoleID, u[0].Gender).
		AddRow(u[1].ID, u[1].Name, u[1].RoleID, u[1].Gender)

	mock.ExpectQuery(query).WillReturnRows(rows)
	//mock.ExpectQuery(query).WithArgs(1).WillReturnError(err)
	res,err := dbHandler.selectData()
	fmt.Println(res)
	assert.Equal(t, u,res)
	//assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestInsert(t *testing.T){
	db,mock,err := sqlmock.New()
	dbHandler := &sqlDb{db}
	defer func() {
		dbHandler.Close()
	}()
	query := "INSERT INTO Employee"

	//rows := sqlmock.NewRows([]string{"name", "roleId", "gender"}).
	//	AddRow(u[0].Name, u[0].RoleID, u[0].Gender).
	//	AddRow(u[1].Name, u[1].RoleID, u[1].Gender)

	mock.ExpectExec(query).WithArgs(u[0].Name, u[0].RoleID, u[0].Gender).WillReturnResult(sqlmock.NewResult(1,1))

	lastID := dbHandler.insertIntoEmployee(u[0])
	assert.Equal(t, u[0].ID,lastID)
	//assert.NotNil(t, res)
	assert.NoError(t, err)
}


func TestUpdate(t *testing.T) {
	db,mock,err := sqlmock.New()
	dbHandler := &sqlDb{db}
	defer func() {
		dbHandler.Close()
	}()
	query := "UPDATE Employee"

	rows := sqlmock.NewRows([]string{"id","name", "roleId", "gender"}).
		AddRow(u[0].ID,u[0].Name, 102, u[0].Gender).
		AddRow(u[1].ID,u[1].Name, u[1].RoleID, u[1].Gender)
    fmt.Println(rows)
	mock.ExpectExec(query).WithArgs(u[0].Name, u[0].RoleID, u[0].Gender,u[0].ID).WillReturnResult(sqlmock.NewResult(0,1))

	rowsAffected := dbHandler.updatEmployee(u[0],u[0].ID)

	assert.NotNil(t, rowsAffected)
	assert.NoError(t, err)
	fmt.Printf("Updated Result:\n")
	TestSelect(t)

}
