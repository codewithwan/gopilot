package domain

import "time"

// URL Shortener models
type ShortURL struct {
	ID          int64      `json:"id"`
	Code        string     `json:"code"`
	OriginalURL string     `json:"original_url"`
	Alias       *string    `json:"alias,omitempty"`
	Clicks      int64      `json:"clicks"`
	IsPublic    bool       `json:"is_public"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateShortURLRequest struct {
	OriginalURL string  `json:"original_url" binding:"required,url"`
	Alias       *string `json:"alias" binding:"omitempty,min=3,max=50,alphanum"`
	ExpireIn    *int    `json:"expire_in" binding:"omitempty,min=1"` // in hours
	IsPublic    *bool   `json:"is_public"`
}

type URLClickLog struct {
	ID         int64     `json:"id"`
	ShortURLID int64     `json:"short_url_id"`
	Referrer   *string   `json:"referrer,omitempty"`
	UserAgent  *string   `json:"user_agent,omitempty"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	ClickedAt  time.Time `json:"clicked_at"`
}

// Pastebin models
type Paste struct {
	ID           string     `json:"id"`
	Title        *string    `json:"title,omitempty"`
	Content      string     `json:"content"`
	Syntax       *string    `json:"syntax,omitempty"`
	IsPublic     bool       `json:"is_public"`
	IsCompressed bool       `json:"is_compressed"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreatePasteRequest struct {
	Title      *string `json:"title" binding:"omitempty,max=255"`
	Content    string  `json:"content" binding:"required"`
	Syntax     *string `json:"syntax" binding:"omitempty,max=50"`
	IsPublic   *bool   `json:"is_public"`
	ExpireIn   *int    `json:"expire_in" binding:"omitempty,min=1"` // in hours
	Compressed *bool   `json:"compressed"`
}

// QR Code models
type QRCode struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Format    string    `json:"format"`
	Size      int       `json:"size"`
	ImageData []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type GenerateQRRequest struct {
	Text   string  `json:"text" binding:"required,max=1000"`
	Size   *int    `json:"size" binding:"omitempty,min=64,max=2048"`
	Format *string `json:"format" binding:"omitempty,oneof=png svg"`
}

// Hash & Encode models
type HashRequest struct {
	Text      string  `json:"text" binding:"required"`
	Algorithm string  `json:"algorithm" binding:"required,oneof=md5 sha1 sha256 sha512 bcrypt"`
	Salt      *string `json:"salt"`
}

type HashResponse struct {
	Hash      string `json:"hash"`
	Algorithm string `json:"algorithm"`
}

type EncodeRequest struct {
	Text      string `json:"text" binding:"required"`
	Operation string `json:"operation" binding:"required,oneof=base64-encode base64-decode url-encode url-decode hex-encode hex-decode"`
}

type EncodeResponse struct {
	Result    string `json:"result"`
	Operation string `json:"operation"`
}

type GeneratePasswordRequest struct {
	Length         *int  `json:"length" binding:"omitempty,min=8,max=128"`
	IncludeUpper   *bool `json:"include_upper"`
	IncludeLower   *bool `json:"include_lower"`
	IncludeNumbers *bool `json:"include_numbers"`
	IncludeSymbols *bool `json:"include_symbols"`
}

type GeneratePasswordResponse struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
}

// Converter models
type ConvertBaseRequest struct {
	Value    string `json:"value" binding:"required"`
	FromBase *int   `json:"from_base" binding:"omitempty,min=2,max=64"`
	ToBase   int    `json:"to_base" binding:"required,min=2,max=64"`
}

type ConvertBaseResponse struct {
	Original string `json:"original"`
	Result   string `json:"result"`
	FromBase int    `json:"from_base"`
	ToBase   int    `json:"to_base"`
}

type ConvertColorRequest struct {
	Value string `json:"value" binding:"required"`
	To    string `json:"to" binding:"required,oneof=rgb hex hsl"`
}

type ConvertColorResponse struct {
	Original string `json:"original"`
	Result   string `json:"result"`
	Format   string `json:"format"`
}

type ConvertTimeRequest struct {
	Value string `json:"value" binding:"required"`
	From  string `json:"from" binding:"required,oneof=unix iso8601"`
	To    string `json:"to" binding:"required,oneof=unix iso8601 human"`
}

type ConvertTimeResponse struct {
	Original string `json:"original"`
	Result   string `json:"result"`
}

type FormatJSONRequest struct {
	JSON   string `json:"json" binding:"required"`
	Minify *bool  `json:"minify"`
	Indent *int   `json:"indent" binding:"omitempty,min=0,max=8"`
}

type FormatJSONResponse struct {
	Result string `json:"result"`
}

// Generator models
type GenerateUUIDRequest struct {
	Version *int `json:"version" binding:"omitempty,oneof=1 4 7"`
	Count   *int `json:"count" binding:"omitempty,min=1,max=100"`
}

type GenerateUUIDResponse struct {
	UUIDs []string `json:"uuids"`
}

type GenerateTokenRequest struct {
	Length *int    `json:"length" binding:"omitempty,min=16,max=256"`
	Prefix *string `json:"prefix" binding:"omitempty,max=20"`
	Suffix *string `json:"suffix" binding:"omitempty,max=20"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

type GenerateLoremRequest struct {
	Type  string `json:"type" binding:"required,oneof=words sentences paragraphs"`
	Count *int   `json:"count" binding:"omitempty,min=1,max=100"`
}

type GenerateLoremResponse struct {
	Text string `json:"text"`
}

type GenerateFakeUserRequest struct {
	Count *int `json:"count" binding:"omitempty,min=1,max=100"`
}

type FakeUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type GenerateFakeUserResponse struct {
	Users []FakeUser `json:"users"`
}

type GenerateRandomNumberRequest struct {
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
	Type  *string  `json:"type" binding:"omitempty,oneof=int float"`
	Count *int     `json:"count" binding:"omitempty,min=1,max=100"`
}

type GenerateRandomNumberResponse struct {
	Numbers []interface{} `json:"numbers"`
}

// Format models
type ConvertYAMLRequest struct {
	Content string `json:"content" binding:"required"`
	To      string `json:"to" binding:"required,oneof=json yaml"`
}

type ConvertYAMLResponse struct {
	Result string `json:"result"`
}

// Crypto models
type AESRequest struct {
	Operation string `json:"operation" binding:"required,oneof=encrypt decrypt"`
	Text      string `json:"text" binding:"required"`
	Key       string `json:"key" binding:"required,min=16,max=32"`
}

type AESResponse struct {
	Result string `json:"result"`
}

type RSAKeypairResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type RSARequest struct {
	Operation string `json:"operation" binding:"required,oneof=encrypt decrypt"`
	Text      string `json:"text" binding:"required"`
	Key       string `json:"key" binding:"required"`
}

type RSAResponse struct {
	Result string `json:"result"`
}

type HMACRequest struct {
	Operation string  `json:"operation" binding:"required,oneof=sign verify"`
	Text      string  `json:"text" binding:"required"`
	Key       string  `json:"key" binding:"required"`
	Signature *string `json:"signature"` // Required for verify
	Algorithm *string `json:"algorithm" binding:"omitempty,oneof=sha256 sha512"`
}

type HMACResponse struct {
	Signature *string `json:"signature,omitempty"`
	Valid     *bool   `json:"valid,omitempty"`
}
