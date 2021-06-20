package service

import (
	"context"
	"crud/models"
	"crud/pkg/repository"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"io"
	"os"
	"time"
)

const (
	salt       = "alkBDij12g3123iuogfd"
	AuthState  = "randomState"
	SigningKey = "ajkhbdoi123@31dj"
	tokenTTL   = 12 * time.Hour
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

type UserService struct {
	repo             repository.User
	googleAuthConfig *oauth2.Config
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
		googleAuthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_SECRET"),
			RedirectURL:  "http://127.0.0.1:8080/auth/callback/google",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *UserService) SignUp(user models.UserInput) error {
	user.Password = generatePasswordHash(user.Password)
	_, err := s.repo.CreateUser(user)
	return err
}

func (s *UserService) SignIn(user models.UserInput) (models.User, error) {
	return s.repo.GetUser(user.Email, generatePasswordHash(user.Password))
}

func (s *UserService) GenerateToken(user models.User) (string, error) {
	claims := &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	}

	tokenConfig := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenConfig.SignedString([]byte(SigningKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) ParseToken(token *jwt.Token) (int, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token claims are not type of jwt.MapClaims")
	}
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("cant parse token")
	}

	return int(id), nil
}

func (s *UserService) GetGoogleAuthUrl() string {
	return s.googleAuthConfig.AuthCodeURL(AuthState)
}

func (s *UserService) GoogleAuthExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.googleAuthConfig.Exchange(ctx, code)
}

func (s *UserService) GetEmailFromGoogleAccessToken(ctx context.Context, accessToken *oauth2.Token) (string, error) {
	type responseJSON struct {
		Email string `json:"email"`
	}

	client := s.googleAuthConfig.Client(ctx, accessToken)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var rJSON responseJSON
	if err = json.Unmarshal(body, &rJSON); err != nil {
		return "", err
	}

	return rJSON.Email, nil
}

func (s *UserService) GetOrCreateUserByEmail(email string) (models.User, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.repo.CreateUser(models.UserInput{Email: email, Password: ""})
		}
		return u, err
	}
	return u, nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
