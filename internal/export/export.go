package export

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/a-h/templ"
	twi "github.com/gavink97/tailwind-inline"
	"github.com/joho/godotenv"
)

type Image struct {
	Path   string
	Name   string
	Inline bool
}

type ImageParams struct {
	Path string
}

func NewCIDImage(params ImageParams) *Image {
	splits := strings.Split(params.Path, "/")
	i := len(splits)
	return &Image{
		Path:   params.Path,
		Name:   splits[i-1],
		Inline: true,
	}
}

func ExportTemplate(c templ.Component) string {
	ctx := context.Background()
	s, err := templ.ToGoHTML(ctx, c)
	if err != nil {
		slog.Info("There was an error parsing the templ component")
		return ""
	}

	path := "./public/css/styles.css"
	styles, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	tpl := twi.Convert(string(s), styles)
	tpl = twi.TransformImgTags(tpl)

	return tpl
}

// add cdn support
// https://mailtrap.io/blog/embedding-images-in-html-email-have-the-rules-changed/
// https://sendgrid.com/en-us/blog/embedding-images-emails-facts
func PrepareImages(tpl, imagetype string) (string, []Image) {
	err := godotenv.Load("/Users/gavinkondrath/projects/templ-mail/.env")
	if err != nil {
		slog.Error(err.Error())
		return "", nil
	}

	clone := strings.Clone(tpl)
	tag := "<img"
	attr := `src="`

	var arr []Image
	for range strings.Count(clone, tag) {
		fi := strings.Index(clone, tag)
		fi = fi + len(tag)

		li := strings.IndexRune(clone[fi:], '>')
		li = fi + li

		fti := strings.Index(clone[fi:li], attr)
		fti = fi + fti + len(attr)

		lti := strings.IndexRune(clone[fti:li], '"')
		lti = fti + lti

		ref := clone[fti:lti]

		splits := strings.Split(ref, "/")
		path := fmt.Sprintf("%s/bin/images/%s", os.Getenv("PROJECT_PATH"), splits[len(splits)-1])
		src := fmt.Sprintf("cid:%s", splits[len(splits)-1])

		if _, err := os.Stat(path); err != nil {
			slog.Error(err.Error())
			return tpl, nil
		}

		switch imagetype {
		case "cid":
			{
				arr = append(arr, *NewCIDImage(ImageParams{Path: path}))
			}
		case "cdn":
			{
				return tpl, nil
			}
		}

		fi = strings.Index(clone, tag)
		clone = fmt.Sprintf("%s%s%s", clone[:fi], strings.Replace(clone[fi:li], tag, "<cln", 1), clone[li:])

		tpl = fmt.Sprintf("%s%s%s", tpl[:fti], strings.Replace(tpl[fti:lti], ref, src, 1), tpl[lti:])
		clone = fmt.Sprintf("%s%s%s", clone[:fti], strings.Replace(clone[fti:lti], ref, src, 1), clone[lti:])
	}

	return tpl, arr
}

func LivePreview(tpl templ.Component) templ.Component {
	str := ExportTemplate(tpl)

	tem, err := template.New("buff").Parse(str)
	if err != nil {
		slog.Error(err.Error())
	}

	tpl2 := templ.FromGoHTML(tem, 0)
	return tpl2
}
