package handlers

import (
	"fmt"
	"net/http"
	"ledgerly/models"
	"ledgerly/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	PettyCashService *services.PettyCashService
	ExpenseService   *services.ExpenseService
	ReportingService *services.ReportingService
	AuthService      *services.AuthService
}

func NewHandler() *Handler {
	return &Handler{
		PettyCashService: &services.PettyCashService{},
		ExpenseService:   &services.ExpenseService{},
		ReportingService: &services.ReportingService{},
		AuthService:      &services.AuthService{},
	}
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"password123"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// BalanceResponse represents petty cash balance
type BalanceResponse struct {
	Balance float64 `json:"balance" example:"5000.00"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if creds.Username == "" || creds.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	token, err := h.AuthService.Login(creds.Username, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// CreatePettyCashTransaction godoc
// @Summary Create petty cash transaction
// @Description Create a credit or debit petty cash transaction
// @Tags Petty Cash
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body models.PettyCashTransaction true "Transaction details"
// @Success 201 {object} models.PettyCashTransaction
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /petty-cash [post]
func (h *Handler) CreatePettyCashTransaction(c *gin.Context) {
	var tx models.PettyCashTransaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	// Graceful RBAC Handling
	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	role := roleVal.(models.UserRole)

	// Check permissions
	allowed := false
	for _, perm := range models.RolePermissions[role] {
		if perm == models.PermissionPettyCashCreate {
			allowed = true
			break
		}
		// Allow creation of Debit transactions if user has expenses.create permission
		if perm == models.PermissionExpensesCreate && tx.Type == models.TransactionTypeDebit {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.PettyCashService.CreateTransaction(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// ListPettyCashTransactions godoc
// @Summary List petty cash transactions
// @Description Get all petty cash transactions
// @Tags Petty Cash
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PettyCashTransaction
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /petty-cash [get]
func (h *Handler) ListPettyCashTransactions(c *gin.Context) {
	transactions, err := h.PettyCashService.ListTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// GetPettyCashBalance godoc
// @Summary Get petty cash balance
// @Description Get current petty cash balance
// @Tags Petty Cash
// @Produce json
// @Security BearerAuth
// @Success 200 {object} BalanceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /petty-cash/balance [get]
func (h *Handler) GetPettyCashBalance(c *gin.Context) {
	balance, err := h.PettyCashService.GetBalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// CreateExpense godoc
// @Summary Create expense
// @Description Create a new expense record
// @Tags Expenses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param expense body models.Expense true "Expense details"
// @Success 201 {object} models.Expense
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /expenses [post]
func (h *Handler) CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Populate UserID from JWT claims
	if userID, exists := c.Get("user_id"); exists {
		expense.UserID = fmt.Sprintf("%d", userID.(uint))
	}

	if err := h.ExpenseService.CreateExpense(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// ListExpenses godoc
// @Summary List expenses
// @Description Get all expenses
// @Tags Expenses
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Expense
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /expenses [get]
func (h *Handler) ListExpenses(c *gin.Context) {
	expenses, err := h.ExpenseService.ListExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// GetExpenseSummary godoc
// @Summary Get expense summary
// @Description Get expense summary with totals by category
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} services.ExpenseSummary
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/expenses-summary [get]
func (h *Handler) GetExpenseSummary(c *gin.Context) {
	summary, err := h.ReportingService.GetExpenseSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// GetPettyCashSummary godoc
// @Summary Get petty cash summary
// @Description Get petty cash summary with credits, debits, and balance
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} services.PettyCashSummary
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/petty-cash-summary [get]
func (h *Handler) GetPettyCashSummary(c *gin.Context) {
	summary, err := h.ReportingService.GetPettyCashSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
