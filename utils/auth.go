package utils

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

type User struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type AppConfig struct {
	Users []User `yaml:"users"`
}

var (
	appUsers []User

	// JWT config
	jwtSecret []byte

	// logout blacklist: store revoked JWT IDs (jti) until expiry
	revokedJTI   = make(map[string]time.Time)
	revokedMutex sync.RWMutex

	// refresh tokens store: token -> username, expiry
	refreshStore = make(map[string]time.Time)
	refreshOwner = make(map[string]string)
	refreshMutex sync.RWMutex
)

// LoadUsersFromYAML loads users from config YAML file
func LoadUsersFromYAML(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	appUsers = cfg.Users
	if len(jwtSecret) == 0 {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-secret-change" // for development only
		}
		jwtSecret = []byte(secret)
	}
	return nil
}

// Authenticate checks if username/password exists in config
func Authenticate(username, password string) bool {
	for _, u := range appUsers {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// CreateToken issues a JWT with subject=username, expiry, and jti
func CreateToken(username string) (string, error) {
	// 1 hour expiry for training purposes
	expiresAt := time.Now().Add(1 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        generateJTI(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// IsTokenValid verifies JWT signature, expiry, and blacklist
func IsTokenValid(tokenString string) bool {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return false
	}
	// check revocation by jti
	if claims.ID != "" {
		revokedMutex.RLock()
		exp, found := revokedJTI[claims.ID]
		revokedMutex.RUnlock()
		if found && time.Now().Before(exp) {
			return false
		}
	}
	return true
}

// InvalidateToken revokes a JWT by its jti until its expiry time
func InvalidateToken(tokenString string) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return
	}
	if claims.ID == "" || claims.ExpiresAt == nil {
		return
	}
	revokedMutex.Lock()
	revokedJTI[claims.ID] = claims.ExpiresAt.Time
	revokedMutex.Unlock()
}

// ExtractBearerToken extracts Bearer token from Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		// Fallback: coba ambil dari cookie agar akses langsung via browser tetap bisa
		if c, err := r.Cookie("access_token"); err == nil && c != nil && c.Value != "" {
			return strings.TrimSpace(c.Value), nil
		}
		return "", errors.New("missing Authorization header")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization header")
	}
	return strings.TrimSpace(parts[1]), nil
}

// Secure moved to utils/middleware.go

// generateJTI creates a simple unique ID for JWT ID claim
func generateJTI() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

// CreateRefreshToken creates a long-lived opaque refresh token
func CreateRefreshToken(username string) (string, error) {
	token := generateJTI()
	exp := time.Now().Add(7 * 24 * time.Hour)
	refreshMutex.Lock()
	refreshStore[token] = exp
	refreshOwner[token] = username
	refreshMutex.Unlock()
	return token, nil
}

// ValidateAndRotateRefresh validates a refresh token and rotates it
// Returns new access token and new refresh token
func ValidateAndRotateRefresh(old string) (string, string, error) {
	refreshMutex.Lock()
	exp, ok := refreshStore[old]
	username := refreshOwner[old]
	if !ok || time.Now().After(exp) {
		refreshMutex.Unlock()
		return "", "", errors.New("invalid refresh token")
	}
	// revoke old
	delete(refreshStore, old)
	delete(refreshOwner, old)
	refreshMutex.Unlock()

	// mint new pair
	access, err := CreateToken(username)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := CreateRefreshToken(username)
	if err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}
