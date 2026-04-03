package audius

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/freedom-music/provider-service/internal/provider/cache"
	"go.uber.org/zap"
)

type Provider struct {
	baseURL string
	appName string
	bearer  string
	client  *http.Client
	cache   *cache.InMemory
	log     *zap.Logger
}

func NewProvider(baseURL, appName, bearer string, cache *cache.InMemory, log *zap.Logger) *Provider {
	return &Provider{
		baseURL: baseURL,
		appName: appName,
		bearer:  bearer,
		client:  &http.Client{Timeout: 8 * time.Second},
		cache:   cache,
		log:     log,
	}
}

func (p *Provider) SearchTracks(q string, limit, offset int) (Paged[Track], error) {
	params := url.Values{
		"query":       {q},
		"limit":       {fmt.Sprint(limit)},
		"offset":      {fmt.Sprint(offset)},
		"sort_method": {"relevant"},
		"app_name":    {p.appName},
	}
	return p.getTrackList("/tracks/search", params, 2*time.Minute)
}

func (p *Provider) Trending(limit, offset int) (Paged[Track], error) {
	params := url.Values{
		"time":     {"week"},
		"limit":    {fmt.Sprint(limit)},
		"offset":   {fmt.Sprint(offset)},
		"app_name": {p.appName},
	}
	return p.getTrackList("/tracks/trending", params, 1*time.Minute)
}

func (p *Provider) Track(id string) (*Track, error) {
	cacheKey := "track:" + id
	if v, ok := p.cache.Get(cacheKey); ok {
		t := v.(Track)
		return &t, nil
	}
	var resp struct {
		Data AudiusTrack `json:"data"`
	}
	err := p.getJSON("/tracks/"+id, nil, &resp)
	if err != nil {
		return nil, err
	}
	t := normalizeTrack(resp.Data)
	if t.StreamURL == "" {
		t.StreamURL = fmt.Sprintf("%s/tracks/%s/stream?app_name=%s", p.baseURL, resp.Data.ID, url.QueryEscape(p.appName))
	}
	p.cache.Set(cacheKey, t, 5*time.Minute)
	return &t, nil
}

func (p *Provider) StreamURL(id string) (string, error) {
	return fmt.Sprintf("%s/tracks/%s/stream?app_name=%s", p.baseURL, id, url.QueryEscape(p.appName)), nil
}

func (p *Provider) Artist(handle string) (*Artist, error) {
	cacheKey := "artist:" + handle
	if v, ok := p.cache.Get(cacheKey); ok {
		a := v.(Artist)
		return &a, nil
	}
	var resp struct {
		Data AudiusUser `json:"data"`
	}
	err := p.getJSON("/users/handle/"+handle, url.Values{"app_name": {p.appName}}, &resp)
	if err != nil {
		return nil, err
	}
	a := normalizeArtist(resp.Data)
	p.cache.Set(cacheKey, a, 10*time.Minute)
	return &a, nil
}

func (p *Provider) ArtistTracks(handle string, limit, offset int) (Paged[Track], error) {
	params := url.Values{
		"limit":    {fmt.Sprint(limit)},
		"offset":   {fmt.Sprint(offset)},
		"app_name": {p.appName},
	}
	path := "/users/handle/" + handle + "/tracks"
	return p.getTrackList(path, params, 3*time.Minute)
}

func (p *Provider) SearchArtists(q string, limit, offset int) (Paged[Artist], error) {
	params := url.Values{
		"query":    {q},
		"limit":    {fmt.Sprint(limit)},
		"offset":   {fmt.Sprint(offset)},
		"app_name": {p.appName},
	}
	var resp struct {
		Data []AudiusUser `json:"data"`
	}
	err := p.getJSON("/users/search", params, &resp)
	if err != nil {
		return Paged[Artist]{Items: []Artist{}, Limit: limit, Offset: offset, HasMore: false}, err
	}
	items := make([]Artist, 0, len(resp.Data))
	for _, u := range resp.Data {
		items = append(items, normalizeArtist(u))
	}
	return Paged[Artist]{Items: items, Limit: limit, Offset: offset, HasMore: len(items) == limit}, nil
}

func (p *Provider) SearchPlaylists(q string, limit, offset int) (Paged[Playlist], error) {
	params := url.Values{
		"query":    {q},
		"limit":    {fmt.Sprint(limit)},
		"offset":   {fmt.Sprint(offset)},
		"app_name": {p.appName},
	}
	var resp struct {
		Data []AudiusPlaylist `json:"data"`
	}
	err := p.getJSON("/playlists/search", params, &resp)
	if err != nil {
		return Paged[Playlist]{Items: []Playlist{}, Limit: limit, Offset: offset, HasMore: false}, err
	}
	items := make([]Playlist, 0, len(resp.Data))
	for _, p := range resp.Data {
		items = append(items, normalizePlaylist(p))
	}
	return Paged[Playlist]{Items: items, Limit: limit, Offset: offset, HasMore: len(items) == limit}, nil
}

func (p *Provider) Playlist(id string) (*Playlist, error) {
	cacheKey := "playlist:" + id
	if v, ok := p.cache.Get(cacheKey); ok {
		pl := v.(Playlist)
		return &pl, nil
	}
	var resp struct {
		Data AudiusPlaylist `json:"data"`
	}
	err := p.getJSON("/playlists/"+id, url.Values{"app_name": {p.appName}}, &resp)
	if err != nil {
		return nil, err
	}
	pl := normalizePlaylist(resp.Data)
	p.cache.Set(cacheKey, pl, 5*time.Minute)
	return &pl, nil
}

func (p *Provider) PlaylistTracks(id string, limit, offset int) (Paged[Track], error) {
	params := url.Values{
		"limit":    {fmt.Sprint(limit)},
		"offset":   {fmt.Sprint(offset)},
		"app_name": {p.appName},
	}
	return p.getTrackList("/playlists/"+id+"/tracks", params, 3*time.Minute)
}

func (p *Provider) getTrackList(path string, params url.Values, ttl time.Duration) (Paged[Track], error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("app_name", p.appName)
	key := path + "?" + params.Encode()
	if v, ok := p.cache.Get(key); ok {
		return v.(Paged[Track]), nil
	}
	var resp struct {
		Data []AudiusTrack `json:"data"`
	}
	err := p.getJSON(path, params, &resp)
	if err != nil {
		return Paged[Track]{Items: []Track{}, Limit: atoi(params.Get("limit")), Offset: atoi(params.Get("offset")), HasMore: false}, err
	}
	items := make([]Track, 0, len(resp.Data))
	for _, t := range resp.Data {
		nt := normalizeTrack(t)
		if nt.StreamURL == "" {
			nt.StreamURL = fmt.Sprintf("%s/tracks/%s/stream?app_name=%s", p.baseURL, t.ID, url.QueryEscape(p.appName))
		}
		items = append(items, nt)
	}
	limit := atoi(params.Get("limit"))
	offset := atoi(params.Get("offset"))
	paged := Paged[Track]{Items: items, Limit: limit, Offset: offset, HasMore: len(items) == limit}
	p.cache.Set(key, paged, ttl)
	return paged, nil
}

func (p *Provider) getJSON(path string, params url.Values, out any) error {
	u := p.baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, _ := http.NewRequest("GET", u, nil)
	if p.bearer != "" {
		req.Header.Set("Authorization", "Bearer "+p.bearer)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		p.log.Warn("audius request failed", zap.String("url", u), zap.Bool("has_bearer", p.bearer != ""), zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		p.log.Warn("audius bad status", zap.String("url", u), zap.Int("status", resp.StatusCode), zap.String("body", string(bodyBytes)), zap.Bool("has_bearer", p.bearer != ""))
		return fmt.Errorf("audius status %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Warn("audius read body failed", zap.String("url", u), zap.Bool("has_bearer", p.bearer != ""), zap.Error(err))
		return err
	}
	if err := json.Unmarshal(bodyBytes, out); err != nil {
		p.log.Warn("audius decode failed", zap.String("url", u), zap.Bool("has_bearer", p.bearer != ""), zap.Error(err), zap.String("body", string(bodyBytes)))
		return err
	}
	return nil
}

func atoi(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}
