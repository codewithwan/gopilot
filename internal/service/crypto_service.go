package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash"
	"io"

	"github.com/codewithwan/gopilot/internal/domain"
)

// CryptoService handles cryptographic operations
type CryptoService struct{}

// NewCryptoService creates a new crypto service
func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

// AESOperation performs AES encryption/decryption
func (s *CryptoService) AESOperation(req *domain.AESRequest) (*domain.AESResponse, error) {
	key := []byte(req.Key)
	
	// Ensure key is proper length (16, 24, or 32 bytes)
	if len(key) < 16 {
		// Pad key if too short
		paddedKey := make([]byte, 16)
		copy(paddedKey, key)
		key = paddedKey
	} else if len(key) > 32 {
		key = key[:32]
	} else if len(key) > 24 {
		key = key[:32]
	} else if len(key) > 16 {
		key = key[:24]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	if req.Operation == "encrypt" {
		plaintext := []byte(req.Text)
		
		// Create GCM mode
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCM: %w", err)
		}

		// Generate nonce
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, fmt.Errorf("failed to generate nonce: %w", err)
		}

		// Encrypt
		ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
		encoded := base64.StdEncoding.EncodeToString(ciphertext)

		return &domain.AESResponse{Result: encoded}, nil
	} else if req.Operation == "decrypt" {
		ciphertext, err := base64.StdEncoding.DecodeString(req.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
		}

		// Create GCM mode
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCM: %w", err)
		}

		nonceSize := gcm.NonceSize()
		if len(ciphertext) < nonceSize {
			return nil, fmt.Errorf("ciphertext too short")
		}

		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt: %w", err)
		}

		return &domain.AESResponse{Result: string(plaintext)}, nil
	}

	return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
}

// GenerateRSAKeypair generates an RSA keypair
func (s *CryptoService) GenerateRSAKeypair() (*domain.RSAKeypairResponse, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Encode private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &domain.RSAKeypairResponse{
		PrivateKey: string(privateKeyPEM),
		PublicKey:  string(publicKeyPEM),
	}, nil
}

// RSAOperation performs RSA encryption/decryption
func (s *CryptoService) RSAOperation(req *domain.RSARequest) (*domain.RSAResponse, error) {
	if req.Operation == "encrypt" {
		// Parse public key
		block, _ := pem.Decode([]byte(req.Key))
		if block == nil {
			return nil, fmt.Errorf("failed to parse PEM block")
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}

		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}

		// Encrypt
		ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, []byte(req.Text), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(ciphertext)
		return &domain.RSAResponse{Result: encoded}, nil
	} else if req.Operation == "decrypt" {
		// Parse private key
		block, _ := pem.Decode([]byte(req.Key))
		if block == nil {
			return nil, fmt.Errorf("failed to parse PEM block")
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		// Decode ciphertext
		ciphertext, err := base64.StdEncoding.DecodeString(req.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
		}

		// Decrypt
		plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt: %w", err)
		}

		return &domain.RSAResponse{Result: string(plaintext)}, nil
	}

	return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
}

// HMACOperation performs HMAC signing/verification
func (s *CryptoService) HMACOperation(req *domain.HMACRequest) (*domain.HMACResponse, error) {
	algorithm := "sha256"
	if req.Algorithm != nil {
		algorithm = *req.Algorithm
	}

	var h func() hash.Hash
	switch algorithm {
	case "sha256":
		h = sha256.New
	case "sha512":
		h = sha512.New
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	if req.Operation == "sign" {
		mac := hmac.New(h, []byte(req.Key))
		mac.Write([]byte(req.Text))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		
		return &domain.HMACResponse{Signature: &signature}, nil
	} else if req.Operation == "verify" {
		if req.Signature == nil {
			return nil, fmt.Errorf("signature is required for verify operation")
		}

		mac := hmac.New(h, []byte(req.Key))
		mac.Write([]byte(req.Text))
		expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		
		valid := hmac.Equal([]byte(expectedSignature), []byte(*req.Signature))
		return &domain.HMACResponse{Valid: &valid}, nil
	}

	return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
}
