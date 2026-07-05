package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

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
	}

	utils.CreatedResponse(ctx, "Table created successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Table updated successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Table status updated successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Table fetched successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Table deleted successfully", nil)
}

func (s *Server) GetAllTablesHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.tableServices.GetAllTablesService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "failed to fetch tables", err)
	}

	utils.PaginatedSuccessResponse(ctx, "Tables fetched successfully", response, *meta)
}
