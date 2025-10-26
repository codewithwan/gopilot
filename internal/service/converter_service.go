package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codewithwan/gopilot/internal/domain"
	"gopkg.in/yaml.v3"
)

// ConverterService handles data conversion operations
type ConverterService struct{}

// NewConverterService creates a new converter service
func NewConverterService() *ConverterService {
	return &ConverterService{}
}

// ConvertBase converts numbers between different bases
func (s *ConverterService) ConvertBase(req *domain.ConvertBaseRequest) (*domain.ConvertBaseResponse, error) {
	fromBase := 10
	if req.FromBase != nil {
		fromBase = *req.FromBase
	}

	// Parse the input value
	num, err := strconv.ParseInt(req.Value, fromBase, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value: %w", err)
	}

	// Convert to target base
	var result string
	switch req.ToBase {
	case 2:
		result = strconv.FormatInt(num, 2)
	case 8:
		result = strconv.FormatInt(num, 8)
	case 10:
		result = strconv.FormatInt(num, 10)
	case 16:
		result = strconv.FormatInt(num, 16)
	default:
		result = strconv.FormatInt(num, req.ToBase)
	}

	return &domain.ConvertBaseResponse{
		Original: req.Value,
		Result:   result,
		FromBase: fromBase,
		ToBase:   req.ToBase,
	}, nil
}

// ConvertColor converts colors between different formats
func (s *ConverterService) ConvertColor(req *domain.ConvertColorRequest) (*domain.ConvertColorResponse, error) {
	value := strings.TrimSpace(req.Value)
	
	// Simple color conversion (basic implementation)
	var result string
	
	if strings.HasPrefix(value, "#") {
		// HEX to RGB or HSL
		if req.To == "rgb" {
			r, g, b, err := s.hexToRGB(value)
			if err != nil {
				return nil, err
			}
			result = fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
		} else {
			return nil, fmt.Errorf("conversion from hex to %s not yet supported", req.To)
		}
	} else if strings.HasPrefix(value, "rgb") {
		// RGB to HEX
		if req.To == "hex" {
			hex, err := s.rgbToHex(value)
			if err != nil {
				return nil, err
			}
			result = hex
		} else {
			return nil, fmt.Errorf("conversion from rgb to %s not yet supported", req.To)
		}
	} else {
		return nil, fmt.Errorf("unsupported color format")
	}

	return &domain.ConvertColorResponse{
		Original: req.Value,
		Result:   result,
		Format:   req.To,
	}, nil
}

// hexToRGB converts hex color to RGB
func (s *ConverterService) hexToRGB(hex string) (int, int, int, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color format")
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %w", err)
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %w", err)
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %w", err)
	}

	// Security: Validate range before conversion
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return 0, 0, 0, fmt.Errorf("color values out of range")
	}

	return int(r), int(g), int(b), nil
}

// rgbToHex converts RGB color to hex
func (s *ConverterService) rgbToHex(rgb string) (string, error) {
	// Parse rgb(r, g, b) format
	rgb = strings.TrimPrefix(rgb, "rgb(")
	rgb = strings.TrimSuffix(rgb, ")")
	parts := strings.Split(rgb, ",")
	
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid RGB format")
	}

	r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))

	return fmt.Sprintf("#%02x%02x%02x", r, g, b), nil
}

// ConvertTime converts time between different formats
func (s *ConverterService) ConvertTime(req *domain.ConvertTimeRequest) (*domain.ConvertTimeResponse, error) {
	var t time.Time
	var err error

	// Parse input
	if req.From == "unix" {
		unix, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse unix timestamp: %w", err)
		}
		t = time.Unix(unix, 0)
	} else if req.From == "iso8601" {
		t, err = time.Parse(time.RFC3339, req.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ISO8601: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported from format: %s", req.From)
	}

	// Convert to output
	var result string
	switch req.To {
	case "unix":
		result = strconv.FormatInt(t.Unix(), 10)
	case "iso8601":
		result = t.Format(time.RFC3339)
	case "human":
		result = t.Format("Mon, 02 Jan 2006 15:04:05 MST")
	default:
		return nil, fmt.Errorf("unsupported to format: %s", req.To)
	}

	return &domain.ConvertTimeResponse{
		Original: req.Value,
		Result:   result,
	}, nil
}

// FormatJSON formats or minifies JSON
func (s *ConverterService) FormatJSON(req *domain.FormatJSONRequest) (*domain.FormatJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(req.JSON), &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	minify := false
	if req.Minify != nil {
		minify = *req.Minify
	}

	var result []byte
	var err error

	if minify {
		result, err = json.Marshal(data)
	} else {
		indent := 2
		if req.Indent != nil {
			indent = *req.Indent
		}
		result, err = json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %w", err)
	}

	return &domain.FormatJSONResponse{
		Result: string(result),
	}, nil
}

// ConvertYAML converts between JSON and YAML
func (s *ConverterService) ConvertYAML(req *domain.ConvertYAMLRequest) (*domain.ConvertYAMLResponse, error) {
	var data interface{}

	if req.To == "json" {
		// YAML to JSON
		if err := yaml.Unmarshal([]byte(req.Content), &data); err != nil {
			return nil, fmt.Errorf("invalid YAML: %w", err)
		}
		result, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to convert to JSON: %w", err)
		}
		return &domain.ConvertYAMLResponse{Result: string(result)}, nil
	} else if req.To == "yaml" {
		// JSON to YAML
		if err := json.Unmarshal([]byte(req.Content), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		result, err := yaml.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to YAML: %w", err)
		}
		return &domain.ConvertYAMLResponse{Result: string(result)}, nil
	}

	return nil, fmt.Errorf("unsupported conversion target: %s", req.To)
}
