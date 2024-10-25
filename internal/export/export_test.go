package export

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestNewCIDImage(t *testing.T) {
	expected := Image{
		Path:   "/public/images/동작절.jpg",
		Name:   "동작절.jpg",
		Inline: true,
	}

	result := NewCIDImage(ImageParams{Path: "/public/images/동작절.jpg"})

	if result.Path != expected.Path || result.Name != expected.Name || result.Inline != expected.Inline {
		t.Errorf("Incorrect Result, result: %v expected: %v", result, expected)
	}

}

func TestPrepareImages(t *testing.T) {
	err := godotenv.Load("/Users/gavinkondrath/projects/templ-mail/.env")
	if err != nil {
		t.Error(err.Error())
	}

	tpl := `<img src="/public/images/동작절.jpg" alt="A temple surrounded by trees"></img>`

	expectedTpl := `<img src="cid:동작절.jpg" alt="A temple surrounded by trees"></img>`

	ppath := os.Getenv("PROJECT_PATH")
	expectedPath := fmt.Sprintf("%s/bin/images/동작절.jpg", ppath)
	expectedName := "동작절.jpg"

	tpl, result := PrepareImages(tpl, "cid")

	if !strings.EqualFold(result[0].Path, expectedPath) || !strings.EqualFold(result[0].Name, expectedName) {
		t.Errorf("Incorrect Result, result path: %s result name: %s expected path: %s expected name: %s", result[0].Path, result[0].Name, expectedPath, expectedName)
	}

	if !strings.EqualFold(tpl, expectedTpl) {
		t.Errorf("Incorrect Result: %s Expected: %s", tpl, expectedTpl)
	}
}
