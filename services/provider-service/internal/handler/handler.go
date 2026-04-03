package handler

import (
	"net/http"
	"strconv"

	"github.com/freedom-music/provider-service/internal/provider/audius"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	a *audius.Provider
}

func New(a *audius.Provider) *Handler { return &Handler{a: a} }

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", httpx.Health)
	api := r.Group("/api/v1/providers/audius")
	{
		api.GET("/search", h.searchTracks)
		api.GET("/tracks/trending", h.trending)
		api.GET("/tracks/:id", h.track)
		api.GET("/tracks/:id/stream", h.stream)
		api.GET("/artists/:handle", h.artist)
		api.GET("/artists/:handle/tracks", h.artistTracks)
		api.GET("/artists/search", h.searchArtists)
		api.GET("/playlists/search", h.searchPlaylists)
		api.GET("/playlists/:id", h.playlist)
		api.GET("/playlists/:id/tracks", h.playlistTracks)
	}
}

func (h *Handler) searchTracks(c *gin.Context) {
	q := c.Query("q")
	limit, offset := parsePaging(c)
	res, err := h.a.SearchTracks(q, limit, offset)
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius search failed")
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) trending(c *gin.Context) {
	limit, offset := parsePaging(c)
	res, err := h.a.Trending(limit, offset)
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius trending failed")
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) track(c *gin.Context) {
	t, err := h.a.Track(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius track failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"track": t})
}

func (h *Handler) stream(c *gin.Context) {
	url, err := h.a.StreamURL(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius stream failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handler) artist(c *gin.Context) {
	a, err := h.a.Artist(c.Param("handle"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius artist failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"artist": a})
}

func (h *Handler) artistTracks(c *gin.Context) {
	limit, offset := parsePaging(c)
	res, err := h.a.ArtistTracks(c.Param("handle"), limit, offset)
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius artist tracks failed")
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) searchArtists(c *gin.Context) {
	q := c.Query("q")
	limit, offset := parsePaging(c)
	res, err := h.a.SearchArtists(q, limit, offset)
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius artists failed")
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) searchPlaylists(c *gin.Context) {
	q := c.Query("q")
	limit, offset := parsePaging(c)
	res, err := h.a.SearchPlaylists(q, limit, offset)
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius playlists failed")
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) playlist(c *gin.Context) {
	pl, err := h.a.Playlist(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius playlist failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlist": pl})
}

func (h *Handler) playlistTracks(c *gin.Context) {
	limit, offset := parsePaging(c)
	res, err := h.a.PlaylistTracks(c.Param("id"), limit, offset)
	if err != nil {
		httpx.JSONError(c, http.StatusBadGateway, sharedErrors.CodeInternal, "audius playlist tracks failed")
		return
	}
	c.JSON(http.StatusOK, res)
}

func parsePaging(c *gin.Context) (int, int) {
	limit := 20
	offset := 0
	if v, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		limit = v
	}
	if v, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil {
		offset = v
	}
	return limit, offset
}
