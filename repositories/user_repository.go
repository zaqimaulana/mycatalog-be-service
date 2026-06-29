package repositories

import (
	"github.com/zaqimaulana/mycatalog-be/config"
	"github.com/zaqimaulana/mycatalog-be/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByFirebaseUID mencari user berdasarkan Firebase UID
func (r *UserRepository) FindByFirebaseUID(uid string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("firebase_uid = ?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

// Create menyimpan user baru ke database
func (r *UserRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// Update memperbarui data user
func (r *UserRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

// FindByID mencari user berdasarkan ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

// UpdateFCMToken update FCM token user
func (r *UserRepository) UpdateFCMToken(userID uint, token string) error {
	return config.DB.Model(&models.User{}).Where("id = ?", userID).Update("fcm_token", token).Error
}
