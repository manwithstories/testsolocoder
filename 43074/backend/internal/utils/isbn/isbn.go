package isbn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"booklibrary/internal/config"
	apperrors "booklibrary/internal/errors"
	"booklibrary/internal/logger"
)

type BookInfo struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"isbn"`
	CoverURL  string `json:"cover_url"`
	Summary   string `json:"summary"`
	Pages     int    `json:"pages"`
}

type openLibraryResponse map[string]struct {
	Title      string   `json:"title"`
	Authors    []struct {
		Name string `json:"name"`
	} `json:"authors"`
	Publishers []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	NumberOfPages int `json:"number_of_pages"`
	Notes         struct {
		Value string `json:"value"`
	} `json:"notes"`
	Cover struct {
		Large string `json:"large"`
		Medium string `json:"medium"`
		Small string `json:"small"`
	} `json:"cover"`
}

type googleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishedDate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			PageCount           int      `json:"pageCount"`
			ImageLinks          struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
				Large          string `json:"large"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

func Validate(isbn string) bool {
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	if len(isbn) == 10 {
		return validateISBN10(isbn)
	} else if len(isbn) == 13 {
		return validateISBN13(isbn)
	}
	return false
}

func validateISBN10(isbn string) bool {
	if len(isbn) != 10 {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		if !isDigit(isbn[i]) {
			return false
		}
		sum += int(isbn[i]-'0') * (10 - i)
	}

	checkChar := isbn[9]
	var checkDigit int
	if checkChar == 'X' || checkChar == 'x' {
		checkDigit = 10
	} else if !isDigit(checkChar) {
		return false
	} else {
		checkDigit = int(checkChar - '0')
	}

	return (sum+checkDigit)%11 == 0
}

func validateISBN13(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}

	sum := 0
	for i := 0; i < 12; i++ {
		if !isDigit(isbn[i]) {
			return false
		}
		digit := int(isbn[i] - '0')
		if i%2 == 1 {
			sum += digit * 3
		} else {
			sum += digit
		}
	}

	if !isDigit(isbn[12]) {
		return false
	}

	checkDigit := int(isbn[12] - '0')
	return (sum+checkDigit)%10 == 0
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func Normalize(isbn string) string {
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")
	return strings.ToUpper(isbn)
}

func FetchBookInfo(isbn string) (*BookInfo, error) {
	normalizedISBN := Normalize(isbn)
	if !Validate(normalizedISBN) {
		return nil, apperrors.ErrISBNInvalid
	}

	client := &http.Client{
		Timeout: time.Duration(config.AppConfig.ISBN.Timeout) * time.Second,
	}

	info, err := fetchFromOpenLibrary(client, normalizedISBN)
	if err == nil && info != nil && info.Title != "" {
		info.ISBN = normalizedISBN
		return info, nil
	}
	logger.Warnf("Failed to fetch from OpenLibrary: %v", err)

	info, err = fetchFromGoogleBooks(client, normalizedISBN)
	if err == nil && info != nil && info.Title != "" {
		info.ISBN = normalizedISBN
		return info, nil
	}
	logger.Warnf("Failed to fetch from Google Books: %v", err)

	return nil, apperrors.ErrISBNNotFound
}

func fetchFromOpenLibrary(client *http.Client, isbn string) (*BookInfo, error) {
	url := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", isbn)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("OpenLibrary response: %s", string(body))

	var result openLibraryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	key := fmt.Sprintf("ISBN:%s", isbn)
	data, ok := result[key]
	if !ok {
		return nil, fmt.Errorf("no data found for ISBN")
	}

	info := &BookInfo{
		Title: data.Title,
		Pages: data.NumberOfPages,
	}

	if len(data.Authors) > 0 {
		info.Author = data.Authors[0].Name
	}

	if len(data.Publishers) > 0 {
		info.Publisher = data.Publishers[0].Name
	}

	if data.Cover.Large != "" {
		info.CoverURL = data.Cover.Large
	} else if data.Cover.Medium != "" {
		info.CoverURL = data.Cover.Medium
	} else if data.Cover.Small != "" {
		info.CoverURL = data.Cover.Small
	}

	info.Summary = data.Notes.Value

	return info, nil
}

func fetchFromGoogleBooks(client *http.Client, isbn string) (*BookInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s", isbn)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Google Books response: %s", string(body))

	var result googleBooksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no data found for ISBN")
	}

	volumeInfo := result.Items[0].VolumeInfo

	info := &BookInfo{
		Title:   volumeInfo.Title,
		Pages:   volumeInfo.PageCount,
		Summary: volumeInfo.Description,
	}

	if len(volumeInfo.Authors) > 0 {
		info.Author = strings.Join(volumeInfo.Authors, ", ")
	}

	info.Publisher = volumeInfo.Publisher

	if volumeInfo.ImageLinks.Large != "" {
		info.CoverURL = volumeInfo.ImageLinks.Large
	} else if volumeInfo.ImageLinks.Thumbnail != "" {
		info.CoverURL = volumeInfo.ImageLinks.Thumbnail
	} else if volumeInfo.ImageLinks.SmallThumbnail != "" {
		info.CoverURL = volumeInfo.ImageLinks.SmallThumbnail
	}

	return info, nil
}

func ExtractPages(description string) int {
	for _, word := range strings.Fields(description) {
		if num, err := strconv.Atoi(strings.TrimRight(word, ".,")); err == nil && num > 10 && num < 5000 {
			return num
		}
	}
	return 0
}
