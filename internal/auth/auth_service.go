package auth

import (
	"context"
	"errors"
	"fmt"
	"service/internal/core"
	"service/internal/repository"
	"service/internal/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *core.UserRequest) (*core.UserResponse, error) {
	// Check if email already exists
	emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Email, nil)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, core.ErrEmailExists
	}

	// Check if phone already exists
	phoneExists, err := s.userRepo.CheckPhoneExists(ctx, req.Phone, nil)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, core.ErrPhoneExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &core.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     req.Role,
		IsActive: true,
	}

	// Set branch ID if provided
	if req.BranchID != nil {
		branchUUID, err := uuid.Parse(*req.BranchID)
		if err != nil {
			return nil, errors.New("invalid branch ID")
		}
		user.BranchID = &branchUUID
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Return user response
	response := user.ToResponse()
	return &response, nil
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, req *core.LoginRequest) (*core.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, core.ErrInvalidPassword
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Return login response
	response := &core.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(24 * time.Hour.Seconds()), // 24 hours
	}

	return response, nil
}

// RefreshToken refreshes access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *core.RefreshTokenRequest) (*core.LoginResponse, error) {
    // Validate refresh token
    claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, core.ErrInvalidToken
	}

    // Check blacklist
    if revoked, err := utils.IsRefreshTokenRevoked(ctx, req.RefreshToken); err == nil && revoked {
        return nil, core.ErrInvalidToken
    }

	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

    // Generate new refresh token (rotation)
    refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

    // Revoke old refresh token
    if parsed, err := utils.ParseRefreshToken(req.RefreshToken); err == nil && parsed.ExpiresAt != nil {
        _ = utils.RevokeRefreshToken(ctx, req.RefreshToken, parsed.ExpiresAt.Time)
    }

	// Return login response
	response := &core.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(24 * time.Hour.Seconds()), // 24 hours
	}

	return response, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
    parsed, err := utils.ParseRefreshToken(refreshToken)
    if err != nil || parsed.ExpiresAt == nil {
        return core.ErrInvalidToken
    }
    return utils.RevokeRefreshToken(ctx, refreshToken, parsed.ExpiresAt.Time)
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *core.ChangePasswordRequest) error {
	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return core.ErrUserNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return core.ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(ctx context.Context, req *core.ForgotPasswordRequest) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate reset token
	resetToken, err := utils.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	// TODO: Send reset email with token
	// For now, just log the token (in production, send via email)
	fmt.Printf("Password reset token for %s: %s\n", user.Email, resetToken)

	return nil
}

// ResetPassword resets user password using reset token
func (s *AuthService) ResetPassword(ctx context.Context, req *core.ResetPasswordRequest) error {
	// Validate reset token
	claims, err := utils.ValidatePasswordResetToken(req.Token)
	if err != nil {
		return core.ErrInvalidToken
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return core.ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// GetProfile retrieves user profile
func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*core.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *core.UserRequest) (*core.UserResponse, error) {
	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	// Check if email already exists (excluding current user)
	emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Email, &userID)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, core.ErrEmailExists
	}

	// Check if phone already exists (excluding current user)
	phoneExists, err := s.userRepo.CheckPhoneExists(ctx, req.Phone, &userID)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, core.ErrPhoneExists
	}

	// Update user fields
	user.Email = req.Email
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Role = req.Role

	// Set branch ID if provided
	if req.BranchID != nil {
		branchUUID, err := uuid.Parse(*req.BranchID)
		if err != nil {
			return nil, errors.New("invalid branch ID")
		}
		user.BranchID = &branchUUID
	}

	// Save user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Return user response
	response := user.ToResponse()
	return &response, nil
}

// ListUsers lists users with pagination and optional filters (admin)
func (s *AuthService) ListUsers(ctx context.Context, page, limit int, role *core.UserRole, branchID *uuid.UUID) ([]*core.User, int64, error) {
	return s.userRepo.List(ctx, page, limit, role, branchID)
}

// GetUser retrieves user by ID (admin)
func (s *AuthService) GetUser(ctx context.Context, id uuid.UUID) (*core.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrUserNotFound
	}
	resp := user.ToResponse()
	return &resp, nil
}

// UpdateUser updates a user by ID (admin)
func (s *AuthService) UpdateUser(ctx context.Context, id uuid.UUID, req *core.UserRequest) (*core.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	// Check uniqueness
	emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Email, &id)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, core.ErrEmailExists
	}

	phoneExists, err := s.userRepo.CheckPhoneExists(ctx, req.Phone, &id)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, core.ErrPhoneExists
	}

	user.Email = req.Email
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Role = req.Role

	if req.BranchID != nil {
		if bid, e := uuid.Parse(*req.BranchID); e == nil {
			user.BranchID = &bid
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

// DeleteUser deletes a user by ID (admin)
func (s *AuthService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// ensure exists
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return core.ErrUserNotFound
	}
	return s.userRepo.Delete(ctx, id)
}

// UpdateFCMToken updates FCM token for a user
func (s *AuthService) UpdateFCMToken(ctx context.Context, userID uuid.UUID, fcmToken string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return core.ErrUserNotFound
	}

	// Update FCM token
	user.FCMToken = fcmToken
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update FCM token: %w", err)
	}

	return nil
}