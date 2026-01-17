package models

type Permission string

const (
	// Auth
	PermissionAuthLogin Permission = "auth.login"

	// Petty Cash
	PermissionPettyCashCreate    Permission = "petty_cash.create" 
	PermissionPettyCashViewList  Permission = "petty_cash.view_list"
	PermissionPettyCashViewBalance Permission = "petty_cash.view_balance"

	// Expenses
	PermissionExpensesCreate  Permission = "expenses.create"
	PermissionExpensesViewOwn Permission = "expenses.view_own"

	// Reports
	PermissionReportsView Permission = "reports.view"
)

var RolePermissions = map[UserRole][]Permission{
	RoleAdmin: {
		PermissionAuthLogin,
		PermissionPettyCashCreate,
		PermissionPettyCashViewList,
		PermissionPettyCashViewBalance,
		PermissionExpensesCreate,
		PermissionExpensesViewOwn,
		PermissionReportsView,
	},
	RoleEmployee: {
		PermissionAuthLogin,
		PermissionPettyCashViewBalance,
		PermissionExpensesCreate,
		PermissionExpensesViewOwn,
	},
}
