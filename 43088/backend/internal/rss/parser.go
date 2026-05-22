package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"podcast-manager/internal/config"
	"podcast-manager/internal/models"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type RSSFeed struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	Channel     Channel  `xml:"channel"`
}

type Channel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	Author        string    `xml:"author"`
	Copyright     string    `xml:"copyright"`
	LastBuildDate string    `xml:"lastBuildDate"`
	PubDate       string    `xml:"pubDate"`
	Image         Image     `xml:"image"`
	ItunesAuthor  string    `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd author"`
	ItunesImage   ItunesImage `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd image"`
	ItunesCategory []ItunesCategory `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd category"`
	Items         []Item    `xml:"item"`
}

type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type ItunesImage struct {
	Href string `xml:"href,attr"`
}

type ItunesCategory struct {
	Text string `xml:"text,attr"`
}

type Item struct {
	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	Description    string    `xml:"description"`
	PubDate        string    `xml:"pubDate"`
	GUID           string    `xml:"guid"`
	Enclosure      Enclosure `xml:"enclosure"`
	ItunesDuration string    `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd duration"`
	ItunesSeason   int       `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd season"`
	ItunesEpisode  int       `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd episode"`
	ItunesEpisodeType string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd episodeType"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

func ParseFeed(feedURL string) (*models.Podcast, []models.Episode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.AppConfig.RSS.Timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		logrus.Errorf("Failed to create request: %v", err)
		return nil, nil, err
	}

	req.Header.Set("User-Agent", config.AppConfig.RSS.UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to fetch RSS feed: %v", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read RSS feed: %v", err)
		return nil, nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		logrus.Errorf("Failed to parse RSS feed: %v", err)
		return nil, nil, err
	}

	podcast := &models.Podcast{
		Title:       feed.Channel.Title,
		Description: feed.Channel.Description,
		FeedURL:     feedURL,
		Website:     feed.Channel.Link,
		Author:      getAuthor(feed.Channel),
		CoverImage:  getCoverImage(feed.Channel),
		Language:    feed.Channel.Language,
		Category:    getCategory(feed.Channel),
		LastChecked: time.Now(),
	}

	var episodes []models.Episode
	for _, item := range feed.Channel.Items {
		episode := models.Episode{
			Title:         item.Title,
			Description:   item.Description,
			GUID:          item.GUID,
			AudioURL:      item.Enclosure.URL,
			AudioType:     item.Enclosure.Type,
			Duration:      parseDuration(item.ItunesDuration),
			PubDate:       parsePubDate(item.PubDate),
			EpisodeType:   item.ItunesEpisodeType,
			SeasonNumber:  item.ItunesSeason,
			EpisodeNumber: item.ItunesEpisode,
			IsNew:         true,
		}
		episodes = append(episodes, episode)
	}

	return podcast, episodes, nil
}

func getAuthor(channel Channel) string {
	if channel.ItunesAuthor != "" {
		return channel.ItunesAuthor
	}
	return channel.Author
}

func getCoverImage(channel Channel) string {
	if channel.ItunesImage.Href != "" {
		return channel.ItunesImage.Href
	}
	return channel.Image.URL
}

func getCategory(channel Channel) string {
	if len(channel.ItunesCategory) > 0 {
		return channel.ItunesCategory[0].Text
	}
	return ""
}

func parseDuration(durationStr string) int {
	if durationStr == "" {
		return 0
	}

	parts := strings.Split(durationStr, ":")
	if len(parts) == 3 {
		hours := parseInt(parts[0])
		minutes := parseInt(parts[1])
		seconds := parseInt(parts[2])
		return hours*3600 + minutes*60 + seconds
	} else if len(parts) == 2 {
		minutes := parseInt(parts[0])
		seconds := parseInt(parts[1])
		return minutes*60 + seconds
	}

	return parseInt(durationStr)
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

func parsePubDate(pubDateStr string) time.Time {
	if pubDateStr == "" {
		return time.Now()
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, pubDateStr); err == nil {
			return t
		}
	}

	logrus.Warnf("Failed to parse date: %s", pubDateStr)
	return time.Now()
}
