package services

import (
	"context"
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testReservationCtx = context.Background()

// MockReservationRepository
type MockReservationRepository struct {
	table        *models.Table
	tables       []models.Table
	reservation  *models.Reservation
	reservations []models.Reservation
	waitlist     *models.Waitlist
	total        int64
	count        int64

	tableErr       error
	reservationErr error
	waitlistErr    error
	createErr      error
	updateErr      error
	countErr       error
}

func newTestGormDB(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectCommit()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	return gormDB
}

func (m *MockReservationRepository) GetTableByIDAndCapacity(_ context.Context, tableID uint, partySize int) (*models.Table, error) {
	return m.table, m.tableErr
}
func (m *MockReservationRepository) GetTablesByCapacity(_ context.Context, partySize int) ([]models.Table, error) {
	return m.tables, m.tableErr
}
func (m *MockReservationRepository) UpdateTableStatus(_ context.Context, _ *gorm.DB, tableID uint, status models.TableStatus) error {
	return m.updateErr
}
func (m *MockReservationRepository) GetReservationsByDateAndSlot(_ context.Context, date string, timeSlot datatypes.Time) ([]models.Reservation, error) {
	return m.reservations, m.reservationErr
}
func (m *MockReservationRepository) CountByTableDateSlot(_ context.Context, tableID uint, date string, timeSlot datatypes.Time) (int64, error) {
	return m.count, m.countErr
}
func (m *MockReservationRepository) Create(_ context.Context, _ *gorm.DB, reservation *models.Reservation) error {
	reservation.ID = 1
	return m.createErr
}
func (m *MockReservationRepository) GetByIDAndUser(_ context.Context, reservationID, userID uint) (*models.Reservation, error) {
	return m.reservation, m.reservationErr
}
func (m *MockReservationRepository) GetByIDWithRelations(_ context.Context, reservationID uint) (*models.Reservation, error) {
	return m.reservation, m.reservationErr
}
func (m *MockReservationRepository) GetByIDAndStatus(_ context.Context, reservationID uint, status models.ReservationStatus) (*models.Reservation, error) {
	return m.reservation, m.reservationErr
}
func (m *MockReservationRepository) UpdateStatus(_ context.Context, _ *gorm.DB, reservationID uint, updates map[string]interface{}) error {
	return m.updateErr
}
func (m *MockReservationRepository) GetAllByUser(_ context.Context, userID uint, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error) {
	return m.reservations, m.total, m.reservationErr
}
func (m *MockReservationRepository) GetAll(_ context.Context, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error) {
	return m.reservations, m.total, m.reservationErr
}
func (m *MockReservationRepository) GetTodayReservations(_ context.Context, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error) {
	return m.reservations, m.total, m.reservationErr
}
func (m *MockReservationRepository) GetFirstWaitlistByDateSlot(_ context.Context, _ *gorm.DB, date interface{}, timeSlot datatypes.Time, partySize int) (*models.Waitlist, error) {
	return m.waitlist, m.waitlistErr
}
func (m *MockReservationRepository) UpdateWaitlistStatus(_ context.Context, _ *gorm.DB, waitlist *models.Waitlist, updates map[string]interface{}) error {
	return m.updateErr
}

func (m *MockReservationRepository) LockTableForUpdate(_ context.Context, _ *gorm.DB, tableID uint) (*models.Table, error) {
	return m.table, m.tableErr
}

func newReservationService(
	repo *MockReservationRepository) *ReservationService {
	return &ReservationService{
		db:              nil,
		reservationRepo: repo,
		redisStore:      redisStore.NewNopCache(),
		publisher:       &MockPublisher{},
	}
}

func newReservationServiceWithDB(t *testing.T, repo *MockReservationRepository) *ReservationService {
	return &ReservationService{
		db:              newTestGormDB(t),
		reservationRepo: repo,
		redisStore:      redisStore.NewNopCache(),
		publisher:       &MockPublisher{},
	}
}

//  CancelReservation Tests

func TestCancelReservation(t *testing.T) {
	userID := uint(1)
	tomorrow := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		reservation *models.Reservation
		repoErr     error
		expectedErr error
	}{
		{
			name: "success",
			reservation: &models.Reservation{
				ID:      1,
				UserID:  userID,
				Status:  models.ReservationStatusConfirmed,
				Date:    tomorrow,
				TableID: 1,
				Table:   models.Table{ID: 1},
				User:    models.User{Email: "test@test.com"},
			},
			expectedErr: nil,
		},
		{
			name:        "not found",
			repoErr:     domain.ErrNotFound,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "already cancelled",
			reservation: &models.Reservation{
				ID:     1,
				UserID: userID,
				Status: models.ReservationStatusCancelled,
				Date:   tomorrow,
			},
			expectedErr: domain.ErrAlreadyCancelled,
		},
		{
			name: "cannot cancel checked in",
			reservation: &models.Reservation{
				ID:     1,
				UserID: userID,
				Status: models.ReservationStatusCheckedIn,
				Date:   tomorrow,
			},
			expectedErr: domain.ErrCannotCancel,
		},
		{
			name: "cannot cancel today",
			reservation: &models.Reservation{
				ID:     1,
				UserID: userID,
				Status: models.ReservationStatusConfirmed,
				Date:   time.Now(), // today → cannot cancel
			},
			expectedErr: domain.ErrCannotCancel,
		},
		{
			name: "wrong user",
			reservation: &models.Reservation{
				ID:     1,
				UserID: 999, // different user
				Status: models.ReservationStatusConfirmed,
				Date:   tomorrow,
			},
			expectedErr: domain.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newReservationServiceWithDB(t, &MockReservationRepository{
				reservation:    tt.reservation,
				reservationErr: tt.repoErr,
				waitlistErr:    domain.ErrNotFound, // no waitlist by default
			})

			response, err := service.CancelReservation(testReservationCtx, userID, 1)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

//  CheckInReservation Tests

func TestCheckInReservation(t *testing.T) {
	tests := []struct {
		name        string
		reservation *models.Reservation
		repoErr     error
		expectedErr error
	}{
		{
			name: "success",
			reservation: &models.Reservation{
				ID:      1,
				UserID:  1,
				Status:  models.ReservationStatusConfirmed,
				TableID: 1,
				User:    models.User{Email: "test@test.com"},
				Table:   models.Table{ID: 1},
			},
			expectedErr: nil,
		},
		{
			name:        "reservation not found",
			repoErr:     domain.ErrNotFound,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "already checked in",
			reservation: &models.Reservation{
				ID:     1,
				Status: models.ReservationStatusCheckedIn,
			},
			expectedErr: domain.ErrAlreadyCheckedIn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newReservationServiceWithDB(t, &MockReservationRepository{
				reservation:    tt.reservation,
				reservationErr: tt.repoErr,
			})

			response, err := service.CheckInReservation(testReservationCtx, 1, 1)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
			}
		})
	}
}

// CheckAvailability Tests

func TestCheckTableAvailability(t *testing.T) {
	tests := []struct {
		name         string
		req          *dto.CheckAvailabilityRequest
		tables       []models.Table
		reservations []models.Reservation
		tableErr     error
		expectedErr  error
	}{
		{
			name: "success with available tables",
			req: &dto.CheckAvailabilityRequest{
				Date:      "2026-07-01",
				TimeSlot:  "20:00:00",
				PartySize: 2,
			},
			tables: []models.Table{
				{ID: 1, Name: "Table 1", Capacity: 4},
				{ID: 2, Name: "Table 2", Capacity: 6},
			},
			reservations: []models.Reservation{}, // none booked
			expectedErr:  nil,
		},
		{
			name: "invalid time slot",
			req: &dto.CheckAvailabilityRequest{
				Date:      "2026-07-01",
				TimeSlot:  "99:00:00",
				PartySize: 2,
			},
			expectedErr: domain.ErrInvalidTimeSlot,
		},
		{
			name: "some tables booked",
			req: &dto.CheckAvailabilityRequest{
				Date:      "2026-07-01",
				TimeSlot:  "20:00:00",
				PartySize: 2,
			},
			tables: []models.Table{
				{ID: 1, Name: "Table 1", Capacity: 4},
				{ID: 2, Name: "Table 2", Capacity: 6},
			},
			reservations: []models.Reservation{
				{TableID: 1}, // table 1 is booked
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newReservationService(&MockReservationRepository{
				tables:       tt.tables,
				reservations: tt.reservations,
				tableErr:     tt.tableErr,
			})

			response, err := service.CheckTableAvailability(testReservationCtx, tt.req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
				assert.Equal(t, tt.req.Date, response.Date)
			}
		})
	}
}

// GetAllUserReservations Tests

func TestGetAllUserReservations_Success(t *testing.T) {
	userID := uint(1)
	service := newReservationService(&MockReservationRepository{
		reservations: []models.Reservation{
			{ID: 1, UserID: userID, Status: models.ReservationStatusConfirmed},
			{ID: 2, UserID: userID, Status: models.ReservationStatusPending},
		},
		total: 2,
	})

	req := &dto.ReservationFilterRequest{Page: 1, PageSize: 10}
	response, err := service.GetAllUserReservations(testReservationCtx, userID, req)

	assert.NoError(t, err)
	assert.Len(t, response.Reservations, 2)
	assert.Equal(t, int64(2), response.Total)
}

func TestCreateReservation_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectCommit()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	repo := &MockReservationRepository{
		table: &models.Table{
			ID:       1,
			Name:     "Table 1",
			Capacity: 4,
		},
		count: 0, // no existing bookings
		reservation: &models.Reservation{
			ID:      1,
			UserID:  1,
			TableID: 1,
			Status:  models.ReservationStatusConfirmed,
			User:    models.User{ID: 1, Email: "test@test.com", FirstName: "Test"},
			Table:   models.Table{ID: 1, Name: "Table 1"},
		},
	}

	service := &ReservationService{
		db:              gormDB,
		reservationRepo: repo,
		redisStore:      redisStore.NewNopCache(),
		publisher:       &MockPublisher{},
	}

	req := &dto.CreateReservationRequest{
		TableID:   1,
		Date:      tomorrow,
		TimeSlot:  "20:00:00",
		PartySize: 2,
	}

	response, err := service.CreateReservation(testReservationCtx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, string(models.ReservationStatusConfirmed), response.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateReservation_TableAlreadyBooked(t *testing.T) {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	service := newReservationServiceWithDB(t, &MockReservationRepository{
		table: &models.Table{ID: 1, Capacity: 4},
		count: 1, // already booked
	})

	req := &dto.CreateReservationRequest{
		TableID:   1,
		Date:      tomorrow,
		TimeSlot:  "20:00:00",
		PartySize: 2,
	}

	response, err := service.CreateReservation(testReservationCtx, req, 1)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrTableAlreadyBooked, err)
}

func TestCreateReservation_PastDate(t *testing.T) {
	service := newReservationService(&MockReservationRepository{})

	req := &dto.CreateReservationRequest{
		TableID:   1,
		Date:      "2020-01-01", // past
		TimeSlot:  "20:00:00",
		PartySize: 2,
	}

	response, err := service.CreateReservation(testReservationCtx, req, 1)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrPastDates, err)
}

func TestCreateReservation_InvalidTimeSlot(t *testing.T) {
	service := newReservationService(&MockReservationRepository{})

	req := &dto.CreateReservationRequest{
		TableID:   1,
		Date:      time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		TimeSlot:  "99:00:00",
		PartySize: 2,
	}

	response, err := service.CreateReservation(testReservationCtx, req, 1)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrInvalidTimeSlot, err)
}
