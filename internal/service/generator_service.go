package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

// GeneratorService handles data generation operations
type GeneratorService struct{}

// NewGeneratorService creates a new generator service
func NewGeneratorService() *GeneratorService {
	return &GeneratorService{}
}

// GenerateUUID generates UUIDs
func (s *GeneratorService) GenerateUUID(req *domain.GenerateUUIDRequest) (*domain.GenerateUUIDResponse, error) {
	version := 4
	if req.Version != nil {
		version = *req.Version
	}

	count := 1
	if req.Count != nil {
		count = *req.Count
	}

	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		var u uuid.UUID
		var err error

		switch version {
		case 1:
			u, err = uuid.NewUUID()
		case 4:
			u, err = uuid.NewRandom()
		case 7:
			u, err = uuid.NewV7()
		default:
			return nil, fmt.Errorf("unsupported UUID version: %d", version)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID: %w", err)
		}

		uuids[i] = u.String()
	}

	return &domain.GenerateUUIDResponse{
		UUIDs: uuids,
	}, nil
}

// GenerateToken generates random tokens
func (s *GeneratorService) GenerateToken(req *domain.GenerateTokenRequest) (*domain.GenerateTokenResponse, error) {
	length := 32
	if req.Length != nil {
		length = *req.Length
	}

	// Generate random bytes
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(b)[:length]

	if req.Prefix != nil && *req.Prefix != "" {
		token = *req.Prefix + "_" + token
	}

	if req.Suffix != nil && *req.Suffix != "" {
		token = token + "_" + *req.Suffix
	}

	return &domain.GenerateTokenResponse{
		Token: token,
	}, nil
}

// GenerateLorem generates lorem ipsum text
func (s *GeneratorService) GenerateLorem(req *domain.GenerateLoremRequest) (*domain.GenerateLoremResponse, error) {
	count := 5
	if req.Count != nil {
		count = *req.Count
	}

	var text string

	switch req.Type {
	case "words":
		words := make([]string, count)
		for i := 0; i < count; i++ {
			words[i] = faker.Word()
		}
		text = strings.Join(words, " ")
	case "sentences":
		sentences := make([]string, count)
		for i := 0; i < count; i++ {
			sentences[i] = faker.Sentence()
		}
		text = strings.Join(sentences, " ")
	case "paragraphs":
		paragraphs := make([]string, count)
		for i := 0; i < count; i++ {
			paragraphs[i] = faker.Paragraph()
		}
		text = strings.Join(paragraphs, "\n\n")
	default:
		return nil, fmt.Errorf("unsupported type: %s", req.Type)
	}

	return &domain.GenerateLoremResponse{
		Text: text,
	}, nil
}

// GenerateFakeUser generates fake user data
func (s *GeneratorService) GenerateFakeUser(req *domain.GenerateFakeUserRequest) (*domain.GenerateFakeUserResponse, error) {
	count := 1
	if req.Count != nil {
		count = *req.Count
	}

	users := make([]domain.FakeUser, count)
	for i := 0; i < count; i++ {
		users[i] = domain.FakeUser{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Username: faker.Username(),
			Phone:    faker.Phonenumber(),
			Address:  fmt.Sprintf("%s, %s", faker.GetRealAddress().Address, faker.GetRealAddress().City),
		}
	}

	return &domain.GenerateFakeUserResponse{
		Users: users,
	}, nil
}

// GenerateRandomNumber generates random numbers
func (s *GeneratorService) GenerateRandomNumber(req *domain.GenerateRandomNumberRequest) (*domain.GenerateRandomNumberResponse, error) {
	min := 0.0
	if req.Min != nil {
		min = *req.Min
	}

	max := 100.0
	if req.Max != nil {
		max = *req.Max
	}

	numType := "int"
	if req.Type != nil {
		numType = *req.Type
	}

	count := 1
	if req.Count != nil {
		count = *req.Count
	}

	numbers := make([]interface{}, count)

	for i := 0; i < count; i++ {
		if numType == "int" {
			// Generate random integer
			diff := int64(max - min)
			n, err := rand.Int(rand.Reader, big.NewInt(diff+1))
			if err != nil {
				return nil, fmt.Errorf("failed to generate random number: %w", err)
			}
			numbers[i] = int(min) + int(n.Int64())
		} else {
			// Generate random float
			n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			if err != nil {
				return nil, fmt.Errorf("failed to generate random number: %w", err)
			}
			fraction := float64(n.Int64()) / float64(math.MaxInt64)
			numbers[i] = min + (max-min)*fraction
		}
	}

	return &domain.GenerateRandomNumberResponse{
		Numbers: numbers,
	}, nil
}
