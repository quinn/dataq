// Code generated by ccf. DO NOT EDIT.
package content

import (
	"embed"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.quinn.io/ccf/content"
)

//go:embed posts
var FS embed.FS
type PostItem = content.ContentItem[Post]

// Initialize loads all content from the embedded filesystem.
// This must be called before using any Get* functions.
func Initialize(e *echo.Echo) error {
	if err := content.LoadItems[Post](FS, "posts"); err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	e.StaticFS("/content", FS)
	return nil
}
// GetPosts returns all posts with their metadata and content.
func GetPosts() ([]PostItem, error) {
	return content.GetItems[Post]()
}
