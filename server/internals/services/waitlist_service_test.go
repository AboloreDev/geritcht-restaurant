package services

import (
	"context"
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

var testWaitlistCtx = context.Background()

// MockWaitlistRepository
type MockWaitlistRepository struct {
	waitlist    *models.Waitlist
	tableCount  int64
	position    int64
	getErr      error
	createErr   error
	countErr    error
	positionErr error
}

func (m *MockWaitlistRepository) CountAvailableTables(_ context.Context, date, timeSlot string, partySize int) (int64, error) {
	return m.tableCount, m.countErr
}
func (m *MockWaitlistRepository) GetByUserDateSlot(_ context.Context, userID uint, date, timeSlot string) (*models.Waitlist, error) {
	return m.waitlist, m.getErr
}
func (m *MockWaitlistRepository) Create(_ context.Context, waitlist *models.Waitlist) error {
	waitlist.ID = 1
	return m.createErr
}
func (m *MockWaitlistRepository) GetPosition(_ context.Context, date string, timeSlot datatypes.Time, createdAt time.Time) (int64, error) {
	return m.position, m.positionErr
}

func newWaitlistService(repo *MockWaitlistRepository) *WaitlistService {
	return NewWaitlistService(repo)
}

// JoinWaitlist Tests
func TestJoinWaitlist(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.JoinWaitlistRequest
		tableCount  int64
		waitlist    *models.Waitlist
		getErr      error
		createErr   error
		expectedErr error
	}{
		{
			name: "success",
			req: &dto.JoinWaitlistRequest{
				Date:      "2026-07-01",
				TimeSlot:  "20:00:00",
				PartySize: 2,
			},
			tableCount:  0,
			getErr:      domain.ErrNotFound,
			expectedErr: nil,
		},
		{
			name: "missing date",
			req: &dto.JoinWaitlistRequest{
				Date:      "",
				TimeSlot:  "20:00:00",
				PartySize: 2,
			},
			expectedErr: domain.ErrInvalidDate,
		},
		{
			name: "invalid party size",
			req: &dto.JoinWaitlistRequest{
				Date:      "2026-07-01",
				TimeSlot:  "20:00:00",
				PartySize: 0,
			},
			expectedErr: domain.ErrInvalidTableCapacity,
		},
		{
			name: "invalid time slot",
			req: &dto.JoinWaitlistRequest{
				Date:      "2026-07-01",
				TimeSlot:  "25:00:00",
				PartySize: 2,
			},
			expectedErr: domain.ErrInvalidTimeSlot,
		},
		{
			name: "table available — no need for waitlist",
			req: &dto.JoinWaitlistRequest{
				Date:      "2026-07-01",
				TimeSlot:  "20:00:00",
				PartySize: 2,
			},
			tableCount:  2,
			expectedErr: domain.ErrTableAvailable,
		},
		{
			name: "already on waitlist",
			req: &dto.JoinWaitlistRequest{
				Date:      "2026-07-01",
				TimeSlot:  "20:00:00",
				PartySize: 2,
			},
			tableCount:  0,
			waitlist:    &models.Waitlist{ID: 1},
			getErr:      nil,
			expectedErr: domain.ErrAlreadyOnWaitlist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newWaitlistService(&MockWaitlistRepository{
				tableCount: tt.tableCount,
				waitlist:   tt.waitlist,
				getErr:     tt.getErr,
				createErr:  tt.createErr,
			})

			response, err := service.JoinWaitlist(testWaitlistCtx, 1, tt.req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
				assert.Equal(t, "20:00:00", response.TimeSlot)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

//  GetWaitlistPosition Tests

func TestGetWaitlistPosition(t *testing.T) {
	tests := []struct {
		name             string
		waitlist         *models.Waitlist
		getErr           error
		position         int64
		expectedPosition int
		expectedErr      error
	}{
		{
			name: "first in queue",
			waitlist: &models.Waitlist{
				ID:        1,
				UserID:    1,
				CreatedAt: time.Now(),
			},
			position:         0, // no one before them
			expectedPosition: 1,
			expectedErr:      nil,
		},
		{
			name: "third in queue",
			waitlist: &models.Waitlist{
				ID:        3,
				UserID:    1,
				CreatedAt: time.Now(),
			},
			position:         2, // two people before them
			expectedPosition: 3,
			expectedErr:      nil,
		},
		{
			name:             "not on waitlist",
			getErr:           domain.ErrNotFound,
			expectedPosition: 0,
			expectedErr:      domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newWaitlistService(&MockWaitlistRepository{
				waitlist: tt.waitlist,
				getErr:   tt.getErr,
				position: tt.position,
			})

			pos, err := service.GetWaitlistPosition(testWaitlistCtx, 1, "2026-07-01", "20:00:00")

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedPosition, pos)
		})
	}
}
