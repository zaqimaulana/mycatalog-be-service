package services

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zaqimaulana/mycatalog-be/config"
	"github.com/zaqimaulana/mycatalog-be/models"
	"github.com/zaqimaulana/mycatalog-be/repositories"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}

// VerifyFirebaseToken verifikasi token dari Firebase,
// pastikan email sudah verified, lalu return Backend JWT
func (s *AuthService) VerifyFirebaseToken(firebaseToken string) (string, *models.User, error) {
	// 1. Verifikasi Firebase ID Token ke server Google
	token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), firebaseToken)
	if err != nil {
		return "", nil, errors.New("firebase token tidak valid atau kadaluarsa")
	}

	// 2. Cek apakah email sudah diverifikasi
	emailVerified, _ := token.Claims["email_verified"].(bool)
	if !emailVerified {
		return "", nil, errors.New("EMAIL_NOT_VERIFIED")
	}

	// 3. Ambil data dari claims Firebase token
	uid := token.UID
	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	// 4. Cari user di database, buat jika belum ada (first time login)
	user, err := s.userRepo.FindByFirebaseUID(uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User pertama kali login — buat user baru
		now := time.Now().Unix()
		user = &models.User{
			FirebaseUID:   uid,
			Email:         email,
			Name:          name,
			Role:          "user",
			EmailVerified: true,
			LastLoginAt:   &now,
		}
		if err := s.userRepo.Create(user); err != nil {
			return "", nil, errors.New("gagal membuat user baru")
		}
	} else if err != nil {
		return "", nil, errors.New("error mengambil data user")
	} else {
		// Update last login
		now := time.Now().Unix()
		user.LastLoginAt = &now
		user.EmailVerified = true
		s.userRepo.Update(user)
	}

	// 5. Generate Backend JWT Token
	jwtToken, err := s.generateJWT(user)
	if err != nil {
		return "", nil, errors.New("gagal membuat token")
	}

	return jwtToken, user, nil
}

// SaveFCMToken simpan FCM device token milik user
func (s *AuthService) SaveFCMToken(userID uint, token string) error {
	return s.userRepo.UpdateFCMToken(userID, token)
}

// generateJWT membuat JWT token dengan payload user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if expireHours == 0 {
		expireHours = 24
	}

	// Claims adalah payload yang disimpan dalam token
	claims := jwt.MapClaims{
		"sub":            user.ID,
		"firebase_uid":   user.FirebaseUID,
		"email":          user.Email,
		"name":           user.Name,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
	}

	// Buat token dengan algoritma HS256 dan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
