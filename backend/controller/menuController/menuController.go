package menuController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/menuService"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService menuService.MenuService
}

func NewMenuController(menuService menuService.MenuService) *MenuController {
	return &MenuController{menuService: menuService}
}

// CreateMenu godoc
// @Summary Create a new menu
// @Description Create a new menu with role assignments
// @Tags menus
// @Accept json
// @Produce json
// @Param menu body model.MenuRequest true "Menu data"
// @Success 201 {object} apiresponse.Response{data=model.Menu}
// @Failure 400 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /menus [post]
// @Security BearerAuth
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var req model.MenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, "INVALID_REQUEST", "Invalid request payload", err, nil)
		return
	}

	menu, err := c.menuService.CreateMenu(&req)
	if err != nil {
		apiresponse.BadRequest(ctx, "CREATE_MENU_FAILED", err.Error(), err, nil)
		return
	}

	apiresponse.Created(ctx, menu, "Menu created successfully", nil)
}

// UpdateMenu godoc
// @Summary Update a menu
// @Description Update an existing menu
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu body model.MenuRequest true "Menu data"
// @Success 200 {object} apiresponse.Response{data=model.Menu}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /menus/{id} [put]
// @Security BearerAuth
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(ctx, "INVALID_MENU_ID", "Invalid menu ID", err, nil)
		return
	}

	var req model.MenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, "INVALID_REQUEST", "Invalid request payload", err, nil)
		return
	}

	menu, err := c.menuService.UpdateMenu(uint(id), &req)
	if err != nil {
		if err.Error() == "menu not found" {
			apiresponse.Error(ctx, http.StatusNotFound, "MENU_NOT_FOUND", "Menu not found", err, nil)
		} else {
			apiresponse.BadRequest(ctx, "UPDATE_MENU_FAILED", err.Error(), err, nil)
		}
		return
	}

	apiresponse.OK(ctx, menu, "Menu updated successfully", nil)
}

// DeleteMenu godoc
// @Summary Delete a menu
// @Description Delete a menu by ID
// @Tags menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} apiresponse.Response
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /menus/{id} [delete]
// @Security BearerAuth
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(ctx, "INVALID_MENU_ID", "Invalid menu ID", err, nil)
		return
	}

	err = c.menuService.DeleteMenu(uint(id))
	if err != nil {
		if err.Error() == "menu not found" {
			apiresponse.Error(ctx, http.StatusNotFound, "MENU_NOT_FOUND", "Menu not found", err, nil)
		} else {
			apiresponse.BadRequest(ctx, "DELETE_MENU_FAILED", err.Error(), err, nil)
		}
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Menu deleted successfully", nil)
}

// GetMenuByID godoc
// @Summary Get menu by ID
// @Description Get a single menu by ID with roles
// @Tags menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} apiresponse.Response{data=model.Menu}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /menus/{id} [get]
// @Security BearerAuth
func (c *MenuController) GetMenuByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(ctx, "INVALID_MENU_ID", "Invalid menu ID", err, nil)
		return
	}

	menu, err := c.menuService.GetMenuByID(uint(id))
	if err != nil {
		if err.Error() == "menu not found" {
			apiresponse.Error(ctx, http.StatusNotFound, "MENU_NOT_FOUND", "Menu not found", err, nil)
		} else {
			apiresponse.InternalServerError(ctx, "GET_MENU_FAILED", "Failed to retrieve menu", err, nil)
		}
		return
	}

	apiresponse.OK(ctx, menu, "Menu retrieved successfully", nil)
}

// GetAllMenus godoc
// @Summary Get all menus
// @Description Get all menus with filters and pagination
// @Tags menus
// @Produce json
// @Param name query string false "Filter by name"
// @Param path query string false "Filter by path"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} apiresponse.Response{data=[]model.Menu}
// @Failure 500 {object} apiresponse.Response
// @Router /menus [get]
// @Security BearerAuth
func (c *MenuController) GetAllMenus(ctx *gin.Context) {
	var filter model.MenuFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apiresponse.BadRequest(ctx, "INVALID_QUERY", "Invalid query parameters", err, nil)
		return
	}

	menus, total, err := c.menuService.GetAllMenus(&filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, "GET_MENUS_FAILED", "Failed to retrieve menus", err, nil)
		return
	}

	// Calculate pagination metadata
	page := filter.Page
	if page == 0 {
		page = 1
	}
	limit := filter.Limit
	if limit == 0 {
		limit = 10
	}

	meta := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
	}

	apiresponse.OK(ctx, menus, "Menus retrieved successfully", meta)
}

// GetUserMenus godoc
// @Summary Get menus for current user
// @Description Get menus accessible by current logged-in user based on their role
// @Tags menus
// @Produce json
// @Success 200 {object} apiresponse.Response{data=[]model.Menu}
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /menus/user [get]
// @Security BearerAuth
func (c *MenuController) GetUserMenus(ctx *gin.Context) {
	// Get role ID from JWT token (set by auth middleware)
	roleID, exists := ctx.Get("role_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Role ID not found in token", nil, nil)
		return
	}

	menus, err := c.menuService.GetUserMenus(roleID.(uint))
	if err != nil {
		apiresponse.InternalServerError(ctx, "GET_USER_MENUS_FAILED", "Failed to retrieve user menus", err, nil)
		return
	}

	apiresponse.OK(ctx, menus, "User menus retrieved successfully", nil)
}

// AssignRolesToMenu godoc
// @Summary Assign roles to menu
// @Description Assign multiple roles to a menu
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param roles body map[string][]uint true "Role IDs" example({"role_ids": [1, 2, 3]})
// @Success 200 {object} apiresponse.Response
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /menus/{id}/roles [post]
// @Security BearerAuth
func (c *MenuController) AssignRolesToMenu(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(ctx, "INVALID_MENU_ID", "Invalid menu ID", err, nil)
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, "INVALID_REQUEST", "Invalid request payload", err, nil)
		return
	}

	err = c.menuService.AssignRolesToMenu(uint(id), req.RoleIDs)
	if err != nil {
		if err.Error() == "menu not found" {
			apiresponse.Error(ctx, http.StatusNotFound, "MENU_NOT_FOUND", "Menu not found", err, nil)
		} else {
			apiresponse.BadRequest(ctx, "ASSIGN_ROLES_FAILED", err.Error(), err, nil)
		}
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Roles assigned successfully", nil)
}
