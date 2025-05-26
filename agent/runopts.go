package agent

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaopandolfi/core"
)

type RunOptions struct {
	Input         string
	StopCondition AgentStopCondition
	Images        []*core.Image
	RunErrs       []error
}

// RunOptionFunc is a function type that modifies RunOptions
type RunOptionFunc func(*RunOptions)

// WithInput sets the input string option
func WithInput(input string) RunOptionFunc {
	return func(opts *RunOptions) {
		opts.Input = input
	}
}

// WithStopCondition sets the stop condition option
func WithStopCondition(stopCondition AgentStopCondition) RunOptionFunc {
	return func(opts *RunOptions) {
		opts.StopCondition = stopCondition
	}
}

func WithImagePath(path string) RunOptionFunc {
	return func(opts *RunOptions) {
		// Read the file
		fileBytes, err := os.ReadFile(path)
		if err != nil {

			panic(err)
		}

		// Get MIME type based on file extension
		mimeType := getMimeType(path)

		// Convert to base64
		encoding := base64.StdEncoding.EncodeToString(fileBytes)

		// Add to Images slice
		opts.Images = append(opts.Images, &core.Image{
			Base64Encoding: encoding,
			MimeType:       mimeType,
		})
	}
}

func WithImageBase64(encoding string, mimeType string) RunOptionFunc {
	return func(opts *RunOptions) {
		opts.Images = append(opts.Images, &core.Image{
			Base64Encoding: encoding,
			MimeType:       mimeType,
		})
	}
}

// Helper function to determine MIME type based on file extension
func getMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream" // Default binary MIME type
	}
}
