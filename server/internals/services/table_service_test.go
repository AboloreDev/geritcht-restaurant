package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/stretchr/testify/assert"
)

var testTableCtx = context.Background()

// MockTableRepository

type MockTableRepository struct {
	table     *models.Table
	tables    []models.Table
	total     int64
	getErr    error
	createErr error
	updateErr error
	deleteErr error
}

func (m *MockTableRepository) GetByName(_ context.Context, name string) (*models.Table, error) {
	return m.table, m.getErr
}
func (m *MockTableRepository) GetByID(_ context.Context, tableID uint) (*models.Table, error) {
	return m.table, m.getErr
}
func (m *MockTableRepository) GetByIDWithRelations(_ context.Context, tableID uint) (*models.Table, error) {
	return m.table, m.getErr
}
func (m *MockTableRepository) Create(_ context.Context, table *models.Table) error {
	table.ID = 1
	return m.createErr
}
func (m *MockTableRepository) Update(_ context.Context, table *models.Table) error {
	return m.updateErr
}
func (m *MockTableRepository) Delete(_ context.Context, tableID uint) error {
	return m.deleteErr
}
func (m *MockTableRepository) GetAll(_ context.Context, page, pageSize int) ([]models.Table, int64, error) {
	return m.tables, m.total, m.getErr
}

func newTableService(repo *MockTableRepository) *TableService {
	return NewTableService(repo, redisStore.NewNopCache())
}

// CreateTable Tests

func TestCreateTable(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.CreateTableRequest
		existing    *models.Table
		getErr      error
		createErr   error
		expectedErr error
	}{
		{
			name:        "success",
			req:         &dto.CreateTableRequest{Name: "Table 1", Capacity: 4},
			getErr:      domain.ErrNotFound, // name not found
			expectedErr: nil,
		},
		{
			name:        "invalid capacity",
			req:         &dto.CreateTableRequest{Name: "Table 1", Capacity: 0},
			expectedErr: domain.ErrInvalidTableCapacity,
		},
		{
			name:        "duplicate name",
			req:         &dto.CreateTableRequest{Name: "Table 1", Capacity: 4},
			existing:    &models.Table{ID: 1, Name: "Table 1"},
			getErr:      nil, // found → conflict
			expectedErr: domain.ErrTableNameConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTableService(&MockTableRepository{
				table:     tt.existing,
				getErr:    tt.getErr,
				createErr: tt.createErr,
			})

			response, err := service.CreateTableService(testTableCtx, tt.req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
				assert.Equal(t, "Table 1", response.Name)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

//  UpdateTable Tests

func TestUpdateTable(t *testing.T) {
	tests := []struct {
		name        string
		table       *models.Table
		getErr      error
		req         *dto.UpdateTableRequest
		expectedErr error
	}{
		{
			name:        "success",
			table:       &models.Table{ID: 1, Name: "Old Name", Capacity: 4},
			req:         &dto.UpdateTableRequest{Name: "New Name"},
			expectedErr: nil,
		},
		{
			name:        "not found",
			getErr:      domain.ErrNotFound,
			req:         &dto.UpdateTableRequest{Name: "New Name"},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTableService(&MockTableRepository{
				table:  tt.table,
				getErr: tt.getErr,
			})

			response, err := service.UpdateTableService(testTableCtx, 1, tt.req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
			}
		})
	}
}

//  UpdateTableStatus Tests

func TestUpdateTableStatus_Success(t *testing.T) {
	service := newTableService(&MockTableRepository{
		table: &models.Table{ID: 1, Status: models.TableStatusAvailable},
	})

	req := &dto.UpdateTableStatusRequest{Status: "occupied"}
	response, err := service.UpdateTableStatusService(testTableCtx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, string(models.TableStatusOccupied), response.Status)
}

func TestUpdateTableStatus_NotFound(t *testing.T) {
	service := newTableService(&MockTableRepository{
		getErr: domain.ErrNotFound,
	})

	req := &dto.UpdateTableStatusRequest{Status: "occupied"}
	response, err := service.UpdateTableStatusService(testTableCtx, 999, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNotFound, err)
}

//  DeleteTable Tests

func TestDeleteTable_Success(t *testing.T) {
	service := newTableService(&MockTableRepository{})

	err := service.DeleteTableService(testTableCtx, 1)

	assert.NoError(t, err)
}

func TestDeleteTable_NotFound(t *testing.T) {
	service := newTableService(&MockTableRepository{
		deleteErr: domain.ErrNotFound,
	})

	err := service.DeleteTableService(testTableCtx, 999)

	assert.Equal(t, domain.ErrNotFound, err)
}

//  GetAllTables Tests

func TestGetAllTables_Success(t *testing.T) {
	service := newTableService(&MockTableRepository{
		tables: []models.Table{
			{ID: 1, Name: "Table 1", Capacity: 4},
			{ID: 2, Name: "Table 2", Capacity: 6},
		},
		total: 2,
	})

	response, meta, err := service.GetAllTablesService(testTableCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
}

func TestGetAllTables_Empty(t *testing.T) {
	service := newTableService(&MockTableRepository{
		tables: []models.Table{},
		total:  0,
	})

	response, meta, err := service.GetAllTablesService(testTableCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 0)
	assert.Equal(t, int64(0), meta.Total)
}
