package auth

import (
	"fmt"
	"net/http"
	"time"

	"explorer451/internal/core"
	"explorer451/internal/models"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// login handles user login via username and password.
func (a *Auth) Login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Validate input format first
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed")
	}

	ctx := c.Request().Context()
	// Get user by username
	user, err := a.repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		a.log.Warn().Msgf("Login attempt for non-existent user: %s", req.Username)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Check if user can login with password
	if !user.PasswordLogin {
		a.log.Warn().Msgf("Password login disabled for user: %s", req.Username)
		return echo.NewHTTPError(http.StatusUnauthorized, "Password login not allowed")
	}

	// Check password
	valid, err := user.CheckPassword(req.Password)
	if err != nil {
		a.log.Error().Err(err).Msgf("Error checking password for user: %s", req.Username)
		return echo.NewHTTPError(http.StatusInternalServerError, "Authentication error")
	}
	if !valid {
		a.log.Warn().Msgf("Invalid password for user: %s", req.Username)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Check if user is active
	if user.Status != "active" {
		a.log.Warn().Msgf("Login attempt for inactive user: %s", req.Username)
		return echo.NewHTTPError(http.StatusForbidden, "Account is not active")
	}

	// Login successful, save session
	if err := a.SaveSession(*user, c); err != nil {
		return err
	}

	a.log.Info().Msgf("User %s (ID: %s) logged in successfully via password.", user.Username, user.ID)
	// Return only non-sensitive user info
	userInfo := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"role":     user.Role,
	}
	return c.JSON(http.StatusOK, userInfo)
}

// logout handles user logout by destroying the session.
func (a *Auth) LogoutHandler(c echo.Context) error {
	if err := a.Logout(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// profile returns the currently authenticated user's profile.
func (a *Auth) Profile(c echo.Context) error {
	// The user object is already set by the auth middleware
	user, ok := c.Get(UserKey).(models.User)
	if !ok {
		// This should technically not happen if middleware is applied correctly
		return echo.NewHTTPError(http.StatusInternalServerError, "User not found in context")
	}

	// Return non-sensitive user info
	userInfo := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"role":     user.Role,
	}
	return c.JSON(http.StatusOK, userInfo)
}

// --- OIDC Handlers (Add these if OIDC is enabled) ---

// oidcLogin redirects the user to the OIDC provider.
func (a *Auth) OIDCLogin(c echo.Context) error {
	if !a.IsOIDCEnabled() {
		return c.String(http.StatusNotImplemented, "OIDC login is not enabled")
	}

	// Generate state and nonce (simple example, enhance for production)
	state, _ := core.GenerateSecureTokenBase64()
	nonce, _ := core.GenerateSecureTokenBase64()

	// Store state/nonce temporarily (e.g., in a short-lived cookie or server-side cache)
	c.SetCookie(&http.Cookie{
		Name:     "oidc_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // Uncomment if using HTTPS, this should be configurable
	})
	c.SetCookie(&http.Cookie{
		Name:     "oidc_nonce",
		Value:    nonce,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // Uncomment if using HTTPS, this should be configurable
	})

	authURL := a.GetOIDCAuthURL(state, nonce)
	return c.Redirect(http.StatusFound, authURL)
}

// oidcCallback handles the callback from the OIDC provider.
func (a *Auth) OIDCCallback(c echo.Context) error {
	if !a.IsOIDCEnabled() {
		return c.String(http.StatusNotImplemented, "OIDC login is not enabled")
	}

	// Retrieve state and nonce from cookies
	stateCookie, err := c.Cookie("oidc_state")
	if err != nil || stateCookie.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "OIDC state cookie missing or empty")
	}
	nonceCookie, err := c.Cookie("oidc_nonce")
	if err != nil || nonceCookie.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "OIDC nonce cookie missing or empty")
	}

	// Clear cookies immediately after retrieving
	stateCookie.MaxAge = -1
	c.SetCookie(stateCookie)
	nonceCookie.MaxAge = -1
	c.SetCookie(nonceCookie)

	// Validate state parameter
	if c.QueryParam("state") != stateCookie.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "OIDC state mismatch")
	}

	// Exchange code for token and verify
	_, claims, err := a.ExchangeOIDCToken(c.QueryParam("code"), nonceCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("OIDC token exchange/verification failed: %v", err))
	}

	// Find user by email claim
	ctx := c.Request().Context()
	user, err := a.repository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Database error finding user by email: %v", err))
	}

	// TODO: Auto-provisioning logic
	if user == nil {
		// User doesn't exist - should we provision a new one or deny login?
		// For now, let's deny login if user doesn't pre-exist
		a.log.Warn().Msgf("OIDC login failed: User with email %s not found.", claims.Email)
		return echo.NewHTTPError(http.StatusForbidden, "User not registered")
	}

	// User exists, check if active
	if user.Status != "active" {
		a.log.Warn().Msgf("OIDC login failed: User account %s is inactive.", claims.Email)
		return echo.NewHTTPError(http.StatusForbidden, "Account is not active")
	}

	// Login successful, save session
	if err := a.SaveSession(*user, c); err != nil {
		return err
	}

	a.log.Info().Msgf("User %s (ID: %s) logged in successfully via OIDC.", user.Username, user.ID)

	// Redirect to the frontend (e.g., dashboard)
	return c.Redirect(http.StatusFound, "/") // Adjust redirect path as needed
}
