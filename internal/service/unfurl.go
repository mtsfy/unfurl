package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type ExtractedData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Site        string `json:"site"`
}

func Fetch(urlStr string) (string, error) {
	if _, err := url.Parse(urlStr); err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if isPopular(urlStr) {
		return usePW(urlStr)
	}

	html, err := useHTTP(urlStr)
	if isSPA(html) || err != nil {
		return usePW(urlStr)
	}

	return html, nil
}

func Extract(data string, baseURL string) (ExtractedData, error) {
	var extracted ExtractedData

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return ExtractedData{}, fmt.Errorf("failed to parse HTML document: %w", err)
	}

	extractors := map[string][]string{
		"title": {
			"meta[property='og:title']",
			"meta[name='twitter:title']",
			"meta[property='twitter:title']",
			"title",
			"meta[itemprop='name']",
			"meta[itemprop='headline']",
			"meta[name='dc.title']",
			"meta[name='DC.title']",
			"meta[name='title']",
			"meta[property='title']",
			"meta[name='page-title']",
			"h1",
			"h2",
			"[data-title]",
		},
		"description": {
			"meta[property='og:description']",
			"meta[name='twitter:description']",
			"meta[property='twitter:description']",
			"meta[name='description']",
			"meta[name='Description']",
			"meta[itemprop='description']",
			"meta[name='dc.description']",
			"meta[name='DC.description']",
			"meta[property='description']",
			"meta[name='page-description']",
			"meta[name='summary']",
			"meta[name='abstract']",
			"meta[name='twitter:summary']",
			"p.description",
			"p.summary",
			"p.lead",
			".description",
		},
		"image": {
			"meta[property='og:image']",
			"meta[property='og:image:url']",
			"meta[property='og:image:secure_url']",
			"meta[name='twitter:image']",
			"meta[name='twitter:image:src']",
			"meta[property='twitter:image']",
			"meta[itemprop='image']",
			"meta[itemprop='thumbnailUrl']",
			"meta[name='image']",
			"meta[property='image']",
			"meta[name='thumbnail']",
			"meta[name='msapplication-TileImage']",
			"link[rel='apple-touch-icon']",
			"link[rel='apple-touch-icon-precomposed']",
			"link[rel='icon'][sizes='192x192']",
			"link[rel='icon'][sizes='180x180']",
			"link[rel='icon'][sizes='32x32']",
			"link[rel='shortcut icon']",
			"link[rel='icon']",
		},
		"site": {
			"meta[property='og:site_name']",
			"meta[name='twitter:site']",
			"meta[property='twitter:site']",
			"meta[name='application-name']",
			"meta[name='apple-mobile-web-app-title']",
			"meta[itemprop='publisher']",
			"meta[name='dc.publisher']",
			"meta[name='DC.publisher']",
			"meta[name='site_name']",
			"meta[name='site-name']",
			"meta[property='site_name']",
			"meta[name='publisher']",
			"meta[name='author']",
			"meta[name='copyright']",
			"meta[name='generator']",
		},
	}

	extracted.Title = trySelectors(doc, extractors["title"])
	extracted.Description = trySelectors(doc, extractors["description"])
	extracted.Image = resolveURL(baseURL, trySelectors(doc, extractors["image"]))
	extracted.Site = trySelectors(doc, extractors["site"])

	if extracted.Site == "" {
		extracted.Site = extracted.Title
	}

	return extracted, nil
}

func trySelectors(doc *goquery.Document, selectors []string) string {
	for _, selector := range selectors {
		if value := getValue(doc, selector); value != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func getValue(doc *goquery.Document, selector string) string {
	elem := doc.Find(selector).First()
	if elem.Length() == 0 {
		return ""
	}

	if href, exists := elem.Attr("href"); exists && elem.Is("link") {
		return href
	}

	if content, exists := elem.Attr("content"); exists {
		return content
	}

	return elem.Text()
}

func resolveURL(baseURL, relativeURL string) string {
	if relativeURL == "" {
		return ""
	}

	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	if strings.HasPrefix(relativeURL, "//") {
		return "https:" + relativeURL
	}

	// Handle absolute paths
	if strings.HasPrefix(relativeURL, "/") {
		// Extract domain from baseURL
		if idx := strings.Index(baseURL, "://"); idx != -1 {
			afterProtocol := baseURL[idx+3:]
			if slashIdx := strings.Index(afterProtocol, "/"); slashIdx != -1 {
				domain := baseURL[:idx+3+slashIdx]
				return domain + relativeURL
			} else {
				return baseURL + relativeURL
			}
		}
	}

	return relativeURL
}

func usePW(urlStr string) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", fmt.Errorf("failed to start playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Timeout: playwright.Float(15000),
	})
	if err != nil {
		return "", fmt.Errorf("failed to launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String("Mozilla/5.0 (compatible; UnfurlBot/1.0)"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create page: %w", err)
	}

	if _, err = page.Goto(urlStr, playwright.PageGotoOptions{
		Timeout:   playwright.Float(15000),
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return "", fmt.Errorf("failed to navigate to %s: %w", urlStr, err)
	}

	html, err := page.Content()
	if err != nil {
		return "", fmt.Errorf("failed to get page content: %w", err)
	}

	return html, nil
}

func useHTTP(urlStr string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; UnfurlBot/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(data), nil
}

func isSPA(html string) bool {
	if strings.Contains(html, "This browser is no longer supported") {
		return true
	}

	signals := []string{
		"data-reactroot", "<div id=\"root\"",
		"data-v-", "<div id=\"app\"",
		"ng-version", "_ngcontent-", "<app-root>",
		"svelte-",
		"__NEXT_DATA__", "<div id=\"__next\">",
		"window.__NUXT__", "<div id=\"__nuxt\">",
		"webpackChunk", "parcelRequire", "vite/dist",
	}

	for _, s := range signals {
		if strings.Contains(html, s) {
			return true
		}
	}

	return len(strings.TrimSpace(html)) < 1000
}

func isPopular(urlStr string) bool {
	sites := []string{
		"x.com", "twitter.com",
		"instagram.com",
		"linkedin.com",
		"tiktok.com",
		"youtube.com",
		"facebook.com",
		"reddit.com",
		"github.com",
		"linkedin.com",
	}

	for _, site := range sites {
		if strings.Contains(urlStr, site) {
			return true
		}
	}
	return false
}

func init() {
	err := playwright.Install() // remove skip install (Dev)
	if err != nil {
		log.Fatal("Failed to install playwright:", err)
	}
}
