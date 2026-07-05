package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

type TableService struct {
	tableRepo  repositories.TableRepositoryInterface
	redisStore interfaces.Cacher
}

func NewTableService(tableRepo repositories.TableRepositoryInterface, redisStore interfaces.Cacher) *TableService {
	return &TableService{tableRepo: tableRepo, redisStore: redisStore}
}

func (s *TableService) CreateTableService(ctx context.Context, req *dto.CreateTableRequest) (*dto.TableResponse, error) {
	if req.Capacity <= 0 {
		return nil, domain.ErrInvalidTableCapacity
	}

	_, err := s.tableRepo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrTableNameConflict
	}

	table := models.Table{
		Name:      req.Name,
		Capacity:  req.Capacity,
		Location:  req.Location,
		Status:    models.TableStatusAvailable,
		CreatedAt: time.Now(),
	}

	if err := s.tableRepo.Create(ctx, &table); err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "table:all:*")
	return mapper.TableResponse(&table), nil
}

func (s *TableService) UpdateTableService(ctx context.Context, tableID uint, req *dto.UpdateTableRequest) (*dto.TableResponse, error) {
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if req.Name != "" {
		table.Name = req.Name
	}
	if req.Location != "" {
		table.Location = req.Location
	}
	if req.Capacity != 0 {
		table.Capacity = req.Capacity
	}

	if err := s.tableRepo.Update(ctx, table); err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "table:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("table:item:%d", tableID))
	return mapper.TableResponse(table), nil
}

func (s *TableService) UpdateTableStatusService(ctx context.Context, tableID uint, req *dto.UpdateTableStatusRequest) (*dto.TableResponse, error) {
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	table.Status = models.TableStatus(req.Status)

	if err := s.tableRepo.Update(ctx, table); err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("table:item:%d", tableID))
	return mapper.TableResponse(table), nil
}

func (s *TableService) GetTableService(ctx context.Context, tableID uint) (*dto.TableDetailResponse, error) {
	table, err := s.tableRepo.GetByIDWithRelations(ctx, tableID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return mapper.TableDetailResponse(table), nil
}

func (s *TableService) DeleteTableService(ctx context.Context, tableID uint) error {
	if err := s.tableRepo.Delete(ctx, tableID); err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "table:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("table:item:%d", tableID))
	return nil
}

func (s *TableService) GetAllTablesService(ctx context.Context, page, pageSize int) ([]*dto.TableResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("table:all:p%d:s%d", page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.TableResponse `json:"data"`
			Meta *utils.PaginatedMeta `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	tables, count, err := s.tableRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.TableResponse, 0, len(tables))
	for _, table := range tables {
		response = append(response, mapper.TableResponse(&table))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{Page: page, Limit: pageSize, Total: count, TotalPages: totalPages}

	cacheData := struct {
		Data []*dto.TableResponse `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}
