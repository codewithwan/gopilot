package service

import (
	"crypto/md5" // #nosec G501 - MD5 is provided as a user-requested feature, not for security
	"crypto/rand"
	"crypto/sha1" // #nosec G505 - SHA1 is provided as a user-requested feature, not for security
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"

	"github.com/codewithwan/gopilot/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// HashService handles hashing and encoding operations
type HashService struct{}

// NewHashService creates a new hash service
func NewHashService() *HashService {
	return &HashService{}
}

// Hash hashes text using the specified algorithm
func (s *HashService) Hash(req *domain.HashRequest) (*domain.HashResponse, error) {
	var hash string
	text := req.Text
	if req.Salt != nil {
		text += *req.Salt
	}

	switch req.Algorithm {
	case "md5":
		h := md5.Sum([]byte(text)) // #nosec G401 - MD5 is provided as a user-requested feature, not for security
		hash = hex.EncodeToString(h[:])
	case "sha1":
		h := sha1.Sum([]byte(text)) // #nosec G401 - SHA1 is provided as a user-requested feature, not for security
		hash = hex.EncodeToString(h[:])
	case "sha256":
		h := sha256.Sum256([]byte(text))
		hash = hex.EncodeToString(h[:])
	case "sha512":
		h := sha512.Sum512([]byte(text))
		hash = hex.EncodeToString(h[:])
	case "bcrypt":
		h, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to generate bcrypt hash: %w", err)
		}
		hash = string(h)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", req.Algorithm)
	}

	return &domain.HashResponse{
		Hash:      hash,
		Algorithm: req.Algorithm,
	}, nil
}

// Encode encodes/decodes text using the specified operation
func (s *HashService) Encode(req *domain.EncodeRequest) (*domain.EncodeResponse, error) {
	var result string
	var err error

	switch req.Operation {
	case "base64-encode":
		result = base64.StdEncoding.EncodeToString([]byte(req.Text))
	case "base64-decode":
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(req.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
		result = string(decoded)
	case "url-encode":
		result = url.QueryEscape(req.Text)
	case "url-decode":
		result, err = url.QueryUnescape(req.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to decode URL: %w", err)
		}
	case "hex-encode":
		result = hex.EncodeToString([]byte(req.Text))
	case "hex-decode":
		var decoded []byte
		decoded, err = hex.DecodeString(req.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex: %w", err)
		}
		result = string(decoded)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
	}

	return &domain.EncodeResponse{
		Result:    result,
		Operation: req.Operation,
	}, nil
}

// GeneratePassword generates a random password
func (s *HashService) GeneratePassword(req *domain.GeneratePasswordRequest) (*domain.GeneratePasswordResponse, error) {
	length := 16
	if req.Length != nil {
		length = *req.Length
		// Security: Limit password length to prevent DoS
		if length > 128 {
			length = 128
		}
	}

	includeUpper := true
	if req.IncludeUpper != nil {
		includeUpper = *req.IncludeUpper
	}

	includeLower := true
	if req.IncludeLower != nil {
		includeLower = *req.IncludeLower
	}

	includeNumbers := true
	if req.IncludeNumbers != nil {
		includeNumbers = *req.IncludeNumbers
	}

	includeSymbols := false
	if req.IncludeSymbols != nil {
		includeSymbols = *req.IncludeSymbols
	}

	charset := ""
	if includeLower {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}
	if includeUpper {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if includeNumbers {
		charset += "0123456789"
	}
	if includeSymbols {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	if charset == "" {
		return nil, fmt.Errorf("at least one character type must be included")
	}

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = charset[num.Int64()]
	}

	return &domain.GeneratePasswordResponse{
		Password: string(password),
		Length:   length,
	}, nil
}
