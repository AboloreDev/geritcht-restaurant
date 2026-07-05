package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateRefreshToken(userID uint, refreshToken string) error {
	token := &models.RefreshToken{
		UserID:    userID,
		TokenHash: utils.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	return r.db.Create(token).Error
}

func (r *AuthRepository) GetValidRefreshToken(ctx context.Context, refreshToken string) (*models.RefreshToken, error) {
	var token models.RefreshToken

	hashedToken := utils.HashToken(refreshToken)

	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND expires_at > ?", hashedToken, time.Now()).
		First(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepository) GetRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken

	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepository) DeleteRefreshToken(token *models.RefreshToken) error {
	return r.db.Delete(token).Error
}

func (r *AuthRepository) DeleteExpiredRefreshTokens(userID uint) error {
	return r.db.Where("user_id = ? AND expires_at < ?", userID, time.Now()).
		Delete(&models.RefreshToken{}).Error
}

func (r *AuthRepository) CreateEmailToken(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) GetValidEmailToken(tokenHash string) (*models.Token, error) {
	var token models.Token

	err := r.db.Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		First(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepository) DeleteEmailToken(token *models.Token) error {
	return r.db.Delete(token).Error
}

func (r *AuthRepository) VerifyUserEmail(user *models.User, token *models.Token) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// update email_verified
		if err := tx.Model(user).
			Update("email_verified", true).Error; err != nil {
			return err
		}

		// delete token
		if err := tx.Delete(token).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *AuthRepository) ResetPassword(user *models.User, token *models.Token, hashedPassword string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(user).Update("password", hashedPassword).Error; err != nil {
			return err
		}

		if err := tx.Delete(token).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", user.ID).
			Delete(&models.RefreshToken{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *AuthRepository) ChangePassword(user *models.User, hashedPassword string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(user).Update("password", hashedPassword).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", user.ID).
			Delete(&models.RefreshToken{}).Error; err != nil {
			return err
		}

		return nil
	})
}
