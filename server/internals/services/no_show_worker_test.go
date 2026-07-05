package services

import (
	"context"
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

var noShowWorkerCtx = context.Background()

type MockReservationNoShowWorkerRepository struct {
	reservation  *models.Reservation
	reservations []models.Reservation

	getErr    error
	updateErr error
}

func (m *MockReservationNoShowWorkerRepository) GetAllReservations(ctx context.Context) ([]models.Reservation, error) {
	return m.reservations, m.getErr
}

func (m *MockReservationNoShowWorkerRepository) MarkReservationNoShow(_ context.Context, reservation *models.Reservation) error {
	return m.updateErr
}

// Helper
func newNoShowRepoWorker(repo *MockReservationNoShowWorkerRepository) *NoShowWorker {
	return NewNoShowWorker(
		&MockPublisher{},
		redisStore.NewNopCache(),
		repo,
	)
}

// Tests
func Test_GetAllReservations_Success(t *testing.T) {
	slotTime := time.Now().Add(-45 * time.Minute)

	timeSlot := datatypes.NewTime(
		slotTime.Hour(),
		slotTime.Minute(),
		slotTime.Second(),
		slotTime.Nanosecond(),
	)
	service := newNoShowRepoWorker(&MockReservationNoShowWorkerRepository{
		reservations: []models.Reservation{
			{
				ID:       1,
				Status:   models.ReservationStatusConfirmed,
				TimeSlot: timeSlot,
				Date:     time.Now(),
			},
			{
				ID:       2,
				Status:   models.ReservationStatusConfirmed,
				TimeSlot: timeSlot,
				Date:     time.Now(),
			},
		},
	})

	response, err := service.noShowWorkerRepo.GetAllReservations(noShowWorkerCtx)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.NotNil(t, response)
}

func Test_GetAllReservations_NotFound(t *testing.T) {
	service := newNoShowRepoWorker(&MockReservationNoShowWorkerRepository{
		getErr: domain.ErrReservationNotFound,
	})

	response, err := service.noShowWorkerRepo.GetAllReservations(noShowWorkerCtx)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrReservationNotFound, err)
	assert.Nil(t, response)
}

func Test_GetAllReservations_Empty(t *testing.T) {
	service := newNoShowRepoWorker(&MockReservationNoShowWorkerRepository{
		reservations: []models.Reservation{},
	})

	response, err := service.noShowWorkerRepo.GetAllReservations(noShowWorkerCtx)
	assert.NoError(t, err)
	assert.Len(t, response, 0)
	assert.NotNil(t, response)
}

func Test_MarkReservationNoShow(t *testing.T) {
	slotTime := time.Now().Add(-45 * time.Minute)

	timeSlot := datatypes.NewTime(
		slotTime.Hour(),
		slotTime.Minute(),
		slotTime.Second(),
		slotTime.Nanosecond(),
	)

	tests := []struct {
		name          string
		reservation   *models.Reservation
		updateErr     error
		getErr        error
		expectedError error
	}{
		{
			name: "successful transaction",
			reservation: &models.Reservation{
				ID:        1,
				Status:    models.ReservationStatusConfirmed,
				Date:      time.Now(),
				TimeSlot:  timeSlot,
				TableID:   1,
				PartySize: 4,
			},
			expectedError: nil,
		},
		{
			name: "reservation update fails",
			reservation: &models.Reservation{
				ID:        1,
				Status:    models.ReservationStatusConfirmed,
				Date:      time.Now(),
				TimeSlot:  timeSlot,
				TableID:   1,
				PartySize: 4,
			},
			updateErr:     domain.ErrInternalServerError,
			expectedError: domain.ErrInternalServerError,
		},
		{
			name: "waitlist lookup fails",
			reservation: &models.Reservation{
				ID:        1,
				Status:    models.ReservationStatusConfirmed,
				Date:      time.Now(),
				TimeSlot:  timeSlot,
				TableID:   1,
				PartySize: 4,
			},
			getErr:        domain.ErrNotFound,
			expectedError: domain.ErrNotFound,
		},
		{
			name: "no waitlist found",
			reservation: &models.Reservation{
				ID:        1,
				Status:    models.ReservationStatusConfirmed,
				Date:      time.Now(),
				TimeSlot:  timeSlot,
				TableID:   1,
				PartySize: 4,
			},
			expectedError: nil,
		},
		{
			name: "waitlist update fails",
			reservation: &models.Reservation{
				ID:        1,
				Status:    models.ReservationStatusConfirmed,
				Date:      time.Now(),
				TimeSlot:  timeSlot,
				TableID:   1,
				PartySize: 4,
			},
			updateErr:     domain.ErrInternalServerError,
			expectedError: domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newNoShowRepoWorker(&MockReservationNoShowWorkerRepository{
				reservation: tt.reservation,
				updateErr:   tt.updateErr,
				getErr:      tt.getErr,
			})

			err := repo.noShowWorkerRepo.MarkReservationNoShow(noShowWorkerCtx, tt.reservation)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
