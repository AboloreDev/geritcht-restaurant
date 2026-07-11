package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create a table
// @Description Create a new restaurant table. Admin access required.
// @Tags Tables
// @Accept json
// @Produce json
// @Param input body dto.CreateTableRequest true "Table details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.TableResponse} "Table created successfully"
// @Failure 400 {object} utils.Response "Invalid request data or invalid table capacity"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 409 {object} utils.Response "Table name already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tables [post]
func (s *Server) CreateTableHandler(ctx *gin.Context) {
	var req dto.CreateTableRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "invalid request body", err)
		return
	}

	response, err := s.tableServices.CreateTableService(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrTableNameConflict:
			utils.ConflictResponse(ctx, "table name already exist", err)
		case domain.ErrInvalidTableCapacity:
			utils.BadRequest(ctx, "invalid table capacity", err)
		default:
			utils.InternalServerError(ctx, "failed to create table", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Table created successfully", response)
}

// @Summary Update a table
// @Description Update an existing restaurant table. Admin access required.
// @Tags Tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Param input body dto.UpdateTableRequest true "Updated table details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.TableResponse} "Table updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or table ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Table not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tables/{id} [patch]
func (s *Server) UpdateTableHandler(ctx *gin.Context) {
	var req dto.UpdateTableRequest
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid table ID", err)
		return
	}
	tableID := uint(id)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "invalid request body", err)
		return
	}

	response, err := s.tableServices.UpdateTableService(ctx.Request.Context(), tableID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "table not found", err)
		default:
			utils.InternalServerError(ctx, "failed to create table", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Table updated successfully", response)
}

// @Summary Update table status
// @Description Update the status of a restaurant table (e.g. Available, Reserved, Occupied). Admin access required.
// @Tags Tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Param input body dto.UpdateTableStatusRequest true "Updated table status"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.TableResponse} "Table status updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or table ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Table not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tables/{id}/status [patch]
func (s *Server) UpdateTableStatusHandler(ctx *gin.Context) {
	var req dto.UpdateTableStatusRequest
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid table ID", err)
		return
	}
	tableID := uint(id)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "invalid request body", err)
		return
	}

	response, err := s.tableServices.UpdateTableStatusService(ctx.Request.Context(), tableID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "table not found", err)
		default:
			utils.InternalServerError(ctx, "failed to update table", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Table status updated successfully", response)
}

// @Summary Get a table
// @Description Retrieve a restaurant table by its ID.
// @Tags Tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Table ID"
// @Success 200 {object} utils.Response{data=dto.TableResponse} "Table retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid table ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Table not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tables/{id} [get]
func (s *Server) GetTableHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid table ID", err)
		return
	}
	tableID := uint(id)

	response, err := s.tableServices.GetTableService(ctx.Request.Context(), tableID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "table not found", err)
		default:
			utils.InternalServerError(ctx, "failed to fetch table", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Table fetched successfully", response)
}

// @Summary Delete a table
// @Description Delete a restaurant table by its ID. Admin access required.
// @Tags Tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Table ID"
// @Success 200 {object} utils.Response "Table deleted successfully"
// @Failure 400 {object} utils.Response "Invalid table ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Table not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tables/{id} [delete]
func (s *Server) DeleteTableHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid table ID", err)
		return
	}
	tableID := uint(id)

	err = s.tableServices.DeleteTableService(ctx.Request.Context(), tableID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "table not found", err)
		default:
			utils.InternalServerError(ctx, "failed to fetch table", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Table deleted successfully", nil)
}

// @Summary List tables
// @Description Retrieve a paginated list of all restaurant tables.
// @Tags Tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.TableResponse} "Tables retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tables [get]
func (s *Server) GetAllTablesHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.tableServices.GetAllTablesService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "failed to fetch tables", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Tables fetched successfully", response, *meta)
}
