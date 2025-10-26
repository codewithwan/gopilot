package handler

import (
	"net/http"
	"strconv"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/middleware"
	"github.com/codewithwan/gopilot/internal/service"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body domain.CreateTodoRequest true "Todo to create"
// @Success 201 {object} domain.Todo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req domain.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todo, err := h.service.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodo godoc
// @Summary Get a todo by ID
// @Description Get a specific todo item by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} domain.Todo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todo, err := h.service.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// ListTodos godoc
// @Summary List all todos
// @Description Get all todos for the authenticated user
// @Tags todos
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/todos [get]
func (h *TodoHandler) ListTodos(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todos, total, err := h.service.List(c.Request.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list todos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos":  todos,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update an existing todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body domain.UpdateTodoRequest true "Todo update data"
// @Success 200 {object} domain.Todo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var req domain.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todo, err := h.service.Update(c.Request.Context(), id, &req, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo item
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
