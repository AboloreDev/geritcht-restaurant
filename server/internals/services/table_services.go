package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type TableService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewTableService(db *gorm.DB, redisStore interfaces.Cacher) *TableService {
	return &TableService{
		redisStore: redisStore,
		db:         db,
	}
}

func (s *TableService) CreateTableService(req *dto.CreateTableRequest) (*dto.TableResponse, error) {
	var table models.Table
	if req.Capacity <= 0 {
		return nil, domain.ErrInvalidTableCapacity
	}

	err := s.db.Where("name = ?", req.Name).First(&table).Error
	if err == nil {
		return nil, domain.ErrTableNameConflict
	}

	table = models.Table{
		Name:      req.Name,
		Capacity:  req.Capacity,
		Location:  req.Location,
		Status:    "available",
		CreatedAt: time.Now(),
	}

	result := s.db.Create(&table)
	if result.Error != nil {
		return nil, result.Error
	}

	s.redisStore.FlushByPattern(ctx, "table:all:*")

	return mapper.TableResponse(&table), nil
}

func (s *TableService) UpdateTableService(tableID uint, req *dto.UpdateTableRequest) (*dto.TableResponse, error) {
	var table models.Table

	err := s.db.Where("id = ?", tableID).First(&table).Error
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

	result := s.db.Save(&table)
	if result.Error != nil {
		return nil, result.Error
	}

	s.redisStore.FlushByPattern(ctx, "table:all:*")
	s.redisStore.Delete(ctx, "table:all")
	s.redisStore.Delete(ctx, fmt.Sprintf("table:item:%d", tableID))

	return mapper.TableResponse(&table), nil
}

func (s *TableService) UpdateTableStatusService(tableID uint, req *dto.UpdateTableStatusRequest) (*dto.TableResponse, error) {
	var table models.Table

	err := s.db.Where("id = ?", tableID).First(&table).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	table.Status = models.TableStatus(req.Status)

	result := s.db.Save(&table)
	if result.Error != nil {
		return nil, result.Error
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("table:item:%d", tableID))

	return mapper.TableResponse(&table), nil
}

func (s *TableService) GetTableService(tableID uint) (*dto.TableDetailResponse, error) {
	var table models.Table

	err := s.db.Preload("Reservations").Preload("Orders").
		Where("id = ? ", tableID).First(&table).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return mapper.TableDetailResponse(&table), nil
}

func (s *TableService) DeleteTableService(tableID uint) error {
	result := s.db.Where("id = ?", tableID).Delete(&models.Table{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	s.redisStore.FlushByPattern(ctx, "table:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("table:item:%d", tableID))

	return nil
}

func (s *TableService) GetAllTablesService(page, pageSize int) ([]*dto.TableResponse, *utils.PaginatedMeta, error) {
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

	offset := utils.Pagination(page, pageSize)

	var tables []models.Table
	var count int64

	s.db.Model(models.Table{}).Count(&count)

	err = s.db.Order("name ASC").Offset(offset).Limit(pageSize).Find(&tables).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.TableResponse, 0, len(tables))

	for _, table := range tables {
		response = append(response, mapper.TableResponse(&table))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.TableResponse `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}

	// store in cache
	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}
