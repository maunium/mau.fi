package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	texttemplate "text/template"
	"time"
	"unsafe"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/gorilla/feeds"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
	"go.mau.fi/util/exerrors"
	"golang.org/x/exp/maps"
	"golang.org/x/net/html"
)

type Post struct {
	Title   string   `yaml:"title"`
	Summary string   `yaml:"summary"`
	Slug    string   `yaml:"slug"`
	Tags    []string `yaml:"tags"`
	Draft   bool     `yaml:"draft"`

	OverrideCreatedAt string `yaml:"override_created_at"`

	HasCodeBlocks bool `yaml:"has_code_blocks"`

	FirstParagraph template.HTML `yaml:"-"`
	Words          int           `yaml:"-"`

	CreatedAt time.Time     `yaml:"-"`
	UpdatedAt time.Time     `yaml:"-"`
	FileName  string        `yaml:"-"`
	Content   template.HTML `yaml:"-"`

	ContentWithoutLinkifiedHeaders string `yaml:"-"`
}

type Tag struct {
	Slug  string
	Name  string
	Posts []*Post

	LastModified time.Time
}

func (p *Post) ToRSS() *feeds.Item {
	return &feeds.Item{
		Title: p.Title,
		Link: &feeds.Link{
			Href: fmt.Sprintf("https://mau.fi/blog/%s/", p.Slug),
		},
		Author: &feeds.Author{
			Name: "Tulir Asokan",
		},
		Description: p.Summary,
		Id:          p.Slug,
		Updated:     p.UpdatedAt,
		Created:     p.CreatedAt,
		Content:     p.ContentWithoutLinkifiedHeaders,
	}
}

var gm = goldmark.New(
	goldmark.WithExtensions(
		extension.Strikethrough,
		extension.Table,
		&frontmatter.Extender{},
		highlighting.NewHighlighting(
			highlighting.WithStyle("solarized-light"),
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
				chromahtml.LineNumbersInTable(true),
				chromahtml.WithAllClasses(true),
				chromahtml.WithLineNumbers(true),
			),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		goldmarkhtml.WithUnsafe(),
	),
)

const PostsDir = "./posts"

var firstParagraphRegex = regexp.MustCompile(`(?s)<p>(.*?)</p>`)
var headerRegex = regexp.MustCompilePOSIX(`<(h[2-6]) id="(.+?)">(.+?)</h[2-6]>`)

func getFileDates(path string) (createdAt, updatedAt time.Time) {
	cmd := exec.Command("git", "log", "--format=%cI", path)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	exerrors.PanicIfNotNil(cmd.Run())
	if stderrBuf.Len() != 0 {
		panic(stderrBuf.String())
	}
	parts := strings.Split(strings.TrimSpace(stdoutBuf.String()), "\n")
	if len(parts) == 1 && parts[0] == "" {
		createdAt = time.Now()
		return
	}
	if len(parts) > 1 {
		updatedAt = exerrors.Must(time.Parse(time.RFC3339, parts[0]))
	}
	createdAt = exerrors.Must(time.Parse(time.RFC3339, parts[len(parts)-1]))
	return
}

func TagNameToSlug(tag string) string {
	return strings.ReplaceAll(strings.ToLower(tag), " ", "-")
}

func WordCount(input string) int {
	return words(exerrors.Must(html.Parse(strings.NewReader(input))))
}

func words(node *html.Node) int {
	switch node.Type {
	case html.TextNode:
		return len(strings.Fields(node.Data))
	case html.ElementNode, html.DocumentNode:
		if node.Data == "pre" {
			// Ignore code blocks entirely
			return 0
		} else if node.Data == "code" {
			// Treat inline code as one word
			return 1
		}
		count := 0
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			count += words(c)
		}
		return count
	default:
		return 0
	}
}

type SitemapParams struct {
	Posts       []*Post
	Tags        []*Tag
	BlogLastMod time.Time
}

var cssURLRegex = regexp.MustCompile(`url\("(/.+?)"\)`)

func inlineifyURLs(input string) string {
	return cssURLRegex.ReplaceAllStringFunc(input, func(s string) string {
		filePath := cssURLRegex.FindStringSubmatch(s)[1]
		var mime string
		switch path.Ext(filePath) {
		case ".svg":
			mime = "image/svg+xml"
		default:
			return s
		}
		data := exerrors.Must(os.ReadFile(".." + filePath))
		encodedData := url.PathEscape(string(data))
		if len(encodedData)+7 < base64.StdEncoding.EncodedLen(len(data)) {
			return fmt.Sprintf(`url("data:%s,%s")`, mime, encodedData)
		} else {
			return fmt.Sprintf(`url("data:%s;base64,%s")`, mime, base64.StdEncoding.EncodeToString(data))
		}
	})
}

func main() {
	sitemapTplData := exerrors.Must(os.ReadFile("../sitemap.xml.tmpl"))
	sitemapTpl := exerrors.Must(texttemplate.New("sitemap").Parse(string(sitemapTplData)))
	tpl := template.New("blog").Funcs(template.FuncMap{
		"include": func(name string) any {
			if strings.HasPrefix(name, "/") {
				name = ".." + name
			}
			data := exerrors.Must(os.ReadFile(name))
			dataStr := unsafe.String(unsafe.SliceData(data), len(data))
			switch filepath.Ext(name) {
			case ".css":
				return template.CSS(inlineifyURLs(dataStr))
			case ".html":
				return template.HTML(dataStr)
			case ".js":
				return template.JS(dataStr)
			default:
				return dataStr
			}
		},
		"safeattr": func(val string) template.HTMLAttr {
			return template.HTMLAttr(val)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"slugifytag":  TagNameToSlug,
		"joinstrings": strings.Join,
	})
	tpl = exerrors.Must(tpl.ParseGlob("templates/*.gohtml"))
	postFiles := exerrors.Must(os.ReadDir(PostsDir))
	postFiles = slices.DeleteFunc(postFiles, func(entry os.DirEntry) bool {
		return !strings.HasSuffix(entry.Name(), ".md")
	})
	posts := make([]*Post, 0, len(postFiles))
	feed := &feeds.Feed{
		Title:       "Tulir Asokan",
		Description: "",
		Link: &feeds.Link{
			Href: "https://mau.fi/blog/",
		},
		Items: make([]*feeds.Item, 0, len(postFiles)),
	}
	feed.Image = &feeds.Image{
		Url:    "https://mau.fi/favicon.png",
		Title:  feed.Title,
		Link:   feed.Link.Href,
		Width:  512,
		Height: 512,
	}
	tags := make(map[string]*Tag)

	for _, file := range postFiles {
		path := filepath.Join(PostsDir, file.Name())
		data := exerrors.Must(os.ReadFile(path))
		var buf strings.Builder
		var meta Post
		ctx := parser.NewContext()
		exerrors.PanicIfNotNil(gm.Convert(data, &buf, parser.WithContext(ctx)))
		exerrors.PanicIfNotNil(frontmatter.Get(ctx).Decode(&meta))
		meta.FileName = "posts/" + file.Name()
		meta.ContentWithoutLinkifiedHeaders = buf.String()
		meta.Content = template.HTML(headerRegex.ReplaceAllString(meta.ContentWithoutLinkifiedHeaders, `<$1 id="$2"><a class="header-anchor" href="#$2">$3</a></$1>`))
		meta.Words = WordCount(meta.ContentWithoutLinkifiedHeaders)
		meta.CreatedAt, meta.UpdatedAt = getFileDates(path)
		if meta.OverrideCreatedAt != "" {
			meta.CreatedAt = exerrors.Must(time.Parse("2006-01-02 15:04:05 -07:00", meta.OverrideCreatedAt))
		}
		meta.HasCodeBlocks = strings.Contains(meta.ContentWithoutLinkifiedHeaders, `class="chroma"`)
		firstParagraphMatch := firstParagraphRegex.FindStringSubmatch(meta.ContentWithoutLinkifiedHeaders)
		if firstParagraphMatch != nil {
			meta.FirstParagraph = template.HTML(firstParagraphMatch[1])
		}

		exerrors.PanicIfNotNil(os.MkdirAll(meta.Slug, 0755))
		mustWriteFile(filepath.Join(meta.Slug, "index.html"), templateExecutor(tpl, "post.gohtml", &meta))

		if meta.Draft {
			continue
		}

		posts = append(posts, &meta)
		feed.Items = append(feed.Items, meta.ToRSS())
		for _, tag := range meta.Tags {
			tagSlug := TagNameToSlug(tag)
			tagMeta, ok := tags[tagSlug]
			if !ok {
				tagMeta = &Tag{
					Slug: tagSlug,
					Name: tag,
				}
				tags[tagSlug] = tagMeta
			}
			tagMeta.Posts = append(tagMeta.Posts, &meta)
			if meta.UpdatedAt.After(tagMeta.LastModified) {
				tagMeta.LastModified = meta.UpdatedAt
			}
			if meta.CreatedAt.After(tagMeta.LastModified) {
				tagMeta.LastModified = meta.CreatedAt
			}
			if meta.UpdatedAt.After(feed.Updated) {
				feed.Updated = meta.UpdatedAt
			}
			if meta.CreatedAt.After(feed.Updated) {
				feed.Updated = meta.CreatedAt
			}
		}
	}
	for _, tag := range tags {
		path := filepath.Join("tags", tag.Slug)
		exerrors.PanicIfNotNil(os.MkdirAll(path, 0755))
		mustWriteFile(filepath.Join(path, "index.html"), templateExecutor(tpl, "index-tag.gohtml", tag))

		tagFeed := &feeds.Feed{
			Title:       fmt.Sprintf("Tulir Asokan - %s", tag.Name),
			Description: "",
			Link: &feeds.Link{
				Href: fmt.Sprintf("https://mau.fi/blog/tags/%s/", tag.Slug),
			},
			Items: make([]*feeds.Item, len(tag.Posts)),
		}
		tagFeed.Image = &feeds.Image{
			Url:    "https://mau.fi/favicon.png",
			Title:  feed.Title,
			Link:   feed.Link.Href,
			Width:  512,
			Height: 512,
		}
		for i, post := range tag.Posts {
			tagFeed.Items[i] = post.ToRSS()
		}
		mustWriteFile(filepath.Join(path, "index.rss"), tagFeed.WriteRss)
		mustWriteFile(filepath.Join(path, "index.atom"), tagFeed.WriteAtom)
		mustWriteFile(filepath.Join(path, "index.json"), tagFeed.WriteJSON)
	}

	mustWriteFile("index.html", templateExecutor(tpl, "index.gohtml", posts))
	mustWriteFile("index.rss", feed.WriteRss)
	mustWriteFile("index.atom", feed.WriteAtom)
	mustWriteFile("index.json", feed.WriteJSON)
	mustWriteFile("../sitemap.xml", func(w io.Writer) error {
		return sitemapTpl.Execute(w, &SitemapParams{
			Posts:       posts,
			Tags:        maps.Values(tags),
			BlogLastMod: feed.Updated,
		})
	})
}

func templateExecutor(tpl *template.Template, name string, data any) func(io.Writer) error {
	return func(w io.Writer) error {
		return tpl.ExecuteTemplate(w, name, data)
	}
}

func mustWriteFile(name string, fn func(io.Writer) error) {
	exerrors.PanicIfNotNil(writeFile(name, fn))
}

func writeFile(name string, fn func(io.Writer) error) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	err = fn(file)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}
