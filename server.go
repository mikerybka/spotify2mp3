package spotify2mp3

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Server struct {
	SongDir             string
	SpotifyClientID     string
	SpotifyClientSecret string
	YoutubeAPIKey       string
}

func (s *Server) dl(spotifyURL string) error {
	cmd := exec.Command("dl-mp3", spotifyURL)
	cmd.Dir = s.SongDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SPOTIFY_CLIENT_ID=%s", s.SpotifyClientID),
		fmt.Sprintf("SPOTIFY_CLIENT_SECRET=%s", s.SpotifyClientSecret),
		fmt.Sprintf("YOUTUBE_API_KEY=%s", s.YoutubeAPIKey),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, out)
	}
	return nil
}

// Save downloads the mp3 to s.SongDir and returns the file reader.
func (s *Server) Save(spotifyURL string) (io.Reader, error) {
	url, err := url.Parse(spotifyURL)
	if err != nil {
		log.Fatal(err)
	}
	spotifyID := strings.TrimPrefix(url.Path, "/track/")
	filename := filepath.Join(s.SongDir, spotifyID+".mp3")

	f, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		err = s.dl(spotifyURL)
		if err != nil {
			return nil, err
		}
		f, err = os.Open(filename)
	}
	return f, err
}
