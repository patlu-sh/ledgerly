package main

import (
	"fmt"
	"ledgerly/db"
	"ledgerly/models"
	"ledgerly/services"
)

func main() {
	db.InitDB()
	authService := &services.AuthService{}

	admin := &models.User{
		Username: "admin",
		Password: "adminpassword",
		Role:     models.RoleAdmin,
	}

	employee := &models.User{
		Username: "employee",
		Password: "employeepassword",
		Role:     models.RoleEmployee,
	}

	if err := authService.Register(admin); err != nil {
		fmt.Println("Error creating admin:", err)
	} else {
		fmt.Println("Admin created")
	}

	if err := authService.Register(employee); err != nil {
		fmt.Println("Error creating employee:", err)
	} else {
		fmt.Println("Employee created")
	}
}
