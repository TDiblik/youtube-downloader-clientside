package main

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

var (
	ErrUnableToParseUrl      = errors.New("unable_to_parse_provided_url")
	ErrNotYoutubeUrl         = errors.New("provided_url_is_not_youtube")
	ErrNoIdProvidedInsideUrl = errors.New("provided_url_does_not_include_id")
)

func GetIdFromYoutubeUrl(query_url string) (string, error) {
	parsed_url, err := url.Parse(query_url)
	if err != nil {
		return "", ErrUnableToParseUrl
	}

	hostname := parsed_url.Hostname()
	hostname = strings.TrimPrefix(hostname, "www.")
	if hostname != "youtube.com" && hostname != "youtu.be" && hostname != "youtube-nocookie.com" {
		return "", ErrNotYoutubeUrl
	}

	// try to get from https://youtu.be/{ID} and https://www.youtube-nocookie.com/**/{id}
	if hostname == "youtu.be" || hostname == "youtube-nocookie.com" {
		possible_id := path.Base(query_url)

		// handles https://youtu.be/ input <-- (possible_id would be youte.be)
		if possible_id == hostname {
			return "", ErrNoIdProvidedInsideUrl
		}

		return possible_id, nil
	}

	// try to get from https://www.youtube.com/watch?v={ID}
	url_query := parsed_url.Query()
	id_from_v_param := url_query.Get("v")
	if !url_query.Has("v") || id_from_v_param == "" {
		return "", ErrNoIdProvidedInsideUrl
	}

	return id_from_v_param, nil
}

// basically the following regex, translated: ^[A-Za-z0-9_-]{11}$
func IsYoutubeIdValid(id string) bool {
	if len(id) != 11 {
		return false
	}

	for _, r := range id {
		if !((r >= 'a' && r <= 'z') || // must be letter between 'a' and 'z'
			(r >= 'A' && r <= 'Z') || // OR must be letter between 'A' and 'Z'
			(r >= '0' && r <= '9') || // OR must be number between '0' and '9'
			r == '_' || r == '-') { // OR must be '-' or '_'
			return false
		}
	}

	return true
}
