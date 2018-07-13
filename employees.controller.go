package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func insertNewEmployee(c *gin.Context) {
	var employee Employees

	if err := c.ShouldBindWith(&employee, binding.JSON); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//check user exist
	err := MySql.Select("id").From("users").
		Where("id", employee.UserId).
		One(&struct{}{})

	if err != nil {
		c.JSON(400, gin.H{"message": "the user_id is wrong!"})
		return
	}

	//check employee existed
	err = MySql.Select("id").From("employees").
		Where("company_id", c.MustGet("company_id").(uint)).
		Where("user_id", employee.UserId).
		One(&struct{}{})

	if err == nil {
		c.JSON(400, gin.H{"message": "employee exist!"})
		return
	}

	//add employee
	employee.CompanyId = c.MustGet("company_id").(uint)
	employee.Status = "pending"
	_, err = MySql.InsertInto("employees").Values(employee).Exec()

	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "successful"})
	return
}

func getAllEmployeesOfCompany(c *gin.Context) {
	type EmployeesWithData struct {
		Id     uint   `db:"id" json:"id"`
		UserId uint   `db:"user_id" json:"user_id"`
		Status string `db:"status" json:"status"`
		Name   string `db:"name" json:"name"`
	}

	var employees []EmployeesWithData

	err := MySql.Select("employees.id", "employees.user_id", "employees.status", "users.name").From("employees").
		Join("users").On("users.id = employees.user_id").
		Where("employees.company_id", c.MustGet("company_id").(uint)).
		All(&employees)

	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, employees)
}

type Employees struct {
	Id        uint   `db:"id" json:"id"`
	CompanyId uint   `db:"company_id" json:"company_id"`
	UserId    uint   `db:"user_id" json:"user_id"`
	Status    string `db:"status" json:"status"`
	Type      string `db:"type" json:"type" binding:"required,oneof=none manager accountant headmaster_accountant technical"`
}