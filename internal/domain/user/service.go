package user

import (
	"encoding/json"
	"fmt"
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/domain/user/dto"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type service struct {
	repo            Repository
	jwtService      auth.JWTService
	googleClientID  string
	superAdminEmail string
}

func NewService(repo Repository, jwtService auth.JWTService, cfg *config.Config) *service {
	return &service{repo: repo, jwtService: jwtService, googleClientID: cfg.GOOGLECLIENTID, superAdminEmail: cfg.SUPERADMINEMAIL}
}

func (s *service) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  RoleUser,
	}
	if req.Email == s.superAdminEmail && s.superAdminEmail != "" {
		user.Role = RoleSuperAdmin
	}

	// hash password and set to user.Password
	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil

}
func (s *service) LoginUser(req dto.LoginRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials // User not found
	}

	// check password
	err = user.checkPassword(req.Password)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}

type googleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Audience      string `json:"aud"`
	EmailVerified bool   `json:"email_verified"`
}

func (s *service) LoginWithGoogle(req dto.GoogleLoginRequest) (*dto.UserResponse, error) {
	if s.googleClientID == "" {
		return nil, fmt.Errorf("google login is not configured")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(req.IDToken))
	if err != nil {
		return nil, fmt.Errorf("verify google token: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, ErrInvalidCredentials
	}

	var info googleTokenInfo
	if err := json.NewDecoder(response.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode google token: %w", err)
	}
	if info.Sub == "" || info.Email == "" || !info.EmailVerified || info.Audience != s.googleClientID {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.GetUserByGoogleID(info.Sub)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = s.repo.GetUserByEmail(info.Email)
		if err != nil {
			return nil, err
		}
	}
	if user == nil {
		user = &User{Name: info.Name, Email: info.Email, GoogleID: info.Sub, Role: RoleUser}
		if err := s.repo.CreateUser(user); err != nil {
			return nil, err
		}
	} else if user.GoogleID == "" {
		user.GoogleID = info.Sub
		if err := s.repo.UpdateUser(user); err != nil {
			return nil, err
		}
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, Token: token, CreatedAt: user.CreatedAt.String()}, nil
}

func (s *service) UpdateRole(id uint, role string) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}
	user.Role = role
	return s.repo.UpdateUser(user)
}
