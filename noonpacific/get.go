package noonpacific

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TODO:fixit
	},
}

// NoonRegexp is the regular expressions for the title of Noon Pacific playlists
var NoonRegexp = regexp.MustCompile(`^NOON \/\/ \d+$`)

// Endpoint is the http API endpoint to get Noon Pacific track data
var Endpoint = "https://api.colormyx.com/v1/noon-pacific/playlists/%d/?detail=true"

// GetPlaylist hits Endpoint to get the playlist with the given ID. If no playlist
// exists, or if the playlist name does not match NoonRegexp, returns an error.
func GetPlaylist(id int) (*Playlist, error) {
	var playlist Playlist
	res, err := client.Get(fmt.Sprintf(Endpoint, id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &playlist)
	if err != nil {
		return nil, err
	}

	if !NoonRegexp.MatchString(playlist.Name) {
		return nil, fmt.Errorf("Invalid playlist name: %v", playlist.Name)
	}

	return &playlist, nil
}