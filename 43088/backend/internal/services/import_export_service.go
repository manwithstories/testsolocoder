package services

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"podcast-manager/internal/database"
	"podcast-manager/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ImportExportService struct{}

func NewImportExportService() *ImportExportService {
	return &ImportExportService{}
}

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title       string `xml:"title"`
	DateCreated string `xml:"dateCreated"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	Text        string    `xml:"text,attr"`
	Title       string    `xml:"title,attr"`
	Type        string    `xml:"type,attr"`
	XMLURL      string    `xml:"xmlUrl,attr"`
	HTMLURL     string    `xml:"htmlUrl,attr"`
	Description string    `xml:"description,attr"`
	Outlines    []Outline `xml:"outline"`
}

func (s *ImportExportService) ExportOPML() (string, error) {
	var podcasts []models.Podcast
	if err := database.DB.Find(&podcasts).Error; err != nil {
		return "", err
	}

	var outlines []Outline
	for _, podcast := range podcasts {
		outlines = append(outlines, Outline{
			Text:        podcast.Title,
			Title:       podcast.Title,
			Type:        "rss",
			XMLURL:      podcast.FeedURL,
			HTMLURL:     podcast.Website,
			Description: podcast.Description,
		})
	}

	opml := OPML{
		Version: "2.0",
		Head: Head{
			Title:       "Podcast Subscriptions",
			DateCreated: time.Now().Format(time.RFC1123Z),
		},
		Body: Body{
			Outlines: outlines,
		},
	}

	output, err := xml.MarshalIndent(opml, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}

func (s *ImportExportService) ImportOPML(content io.Reader) (int, error) {
	body, err := io.ReadAll(content)
	if err != nil {
		return 0, err
	}

	var opml OPML
	if err := xml.Unmarshal(body, &opml); err != nil {
		return 0, err
	}

	podcastService := NewPodcastService()
	importedCount := 0

	var processOutlines func(outlines []Outline)
	processOutlines = func(outlines []Outline) {
		for _, outline := range outlines {
			if outline.XMLURL != "" {
				_, err := podcastService.AddPodcast(outline.XMLURL)
				if err == nil {
					importedCount++
				}
			}
			if len(outline.Outlines) > 0 {
				processOutlines(outline.Outlines)
			}
		}
	}

	processOutlines(opml.Body.Outlines)

	return importedCount, nil
}

func (s *ImportExportService) ExportHistoryCSV() ([][]string, error) {
	var histories []models.ListeningHistory
	if err := database.DB.Preload("Episode.Podcast").Order("start_time DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	records := [][]string{
		{"Date", "Podcast", "Episode", "Duration (seconds)", "Completion (%)", "Start Time", "End Time"},
	}

	for _, h := range histories {
		records = append(records, []string{
			h.StartTime.Format("2006-01-02"),
			h.Episode.Podcast.Title,
			h.Episode.Title,
			strconv.Itoa(h.Duration),
			fmt.Sprintf("%.1f", h.Completion*100),
			h.StartTime.Format(time.RFC3339),
			h.EndTime.Format(time.RFC3339),
		})
	}

	return records, nil
}

func (s *ImportExportService) ExportNotesCSV() ([][]string, error) {
	var notes []models.Note
	if err := database.DB.Preload("Episode.Podcast").Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}

	records := [][]string{
		{"Created At", "Podcast", "Episode", "Timestamp (seconds)", "Content", "Tags"},
	}

	for _, n := range notes {
		tags := ""
		for i, t := range n.Tags {
			if i > 0 {
				tags += "; "
			}
			tags += t
		}

		records = append(records, []string{
			n.CreatedAt.Format(time.RFC3339),
			n.Episode.Podcast.Title,
			n.Episode.Title,
			fmt.Sprintf("%.1f", n.Timestamp),
			n.Content,
			tags,
		})
	}

	return records, nil
}

func (s *ImportExportService) ExportNotesMarkdown(episodeID uuid.UUID) (string, error) {
	var notes []models.Note
	query := database.DB.Preload("Episode.Podcast")
	if episodeID != uuid.Nil {
		query = query.Where("episode_id = ?", episodeID)
	}
	if err := query.Order("timestamp ASC").Find(&notes).Error; err != nil {
		return "", err
	}

	if len(notes) == 0 {
		return "", nil
	}

	content := "# Podcast Notes\n\n"
	content += fmt.Sprintf("Exported: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	currentEpisode := uuid.Nil
	for _, note := range notes {
		if note.EpisodeID != currentEpisode {
			currentEpisode = note.EpisodeID
			content += fmt.Sprintf("## %s - %s\n\n", note.Episode.Podcast.Title, note.Episode.Title)
		}

		minutes := int(note.Timestamp / 60)
		seconds := int(note.Timestamp) % 60
		content += fmt.Sprintf("### [%02d:%02d]\n\n", minutes, seconds)
		content += fmt.Sprintf("%s\n\n", note.Content)

		if len(note.Tags) > 0 {
			content += "Tags: "
			for i, tag := range note.Tags {
				if i > 0 {
					content += ", "
				}
				content += fmt.Sprintf("`%s`", tag)
			}
			content += "\n\n"
		}
	}

	return content, nil
}

func WriteCSV(records [][]string, w io.Writer) error {
	writer := csv.NewWriter(w)
	writer.UseCRLF = true
	return writer.WriteAll(records)
}
