package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authProxy     *httputil.ReverseProxy
	userProxy     *httputil.ReverseProxy
	catalogProxy  *httputil.ReverseProxy
	playlistProxy *httputil.ReverseProxy
	libraryProxy  *httputil.ReverseProxy
	playbackProxy *httputil.ReverseProxy
	mediaProxy    *httputil.ReverseProxy
	eventsProxy   *httputil.ReverseProxy
	searchProxy   *httputil.ReverseProxy
	recoProxy     *httputil.ReverseProxy
	providerProxy *httputil.ReverseProxy
}

func NewRouter(authURL, userURL, catalogURL, playlistURL, libraryURL, playbackURL, mediaURL, eventsURL, searchURL, recoURL, providerURL string) (*Router, error) {
	auth, err := proxyFor(authURL, "/api/v1/auth")
	if err != nil {
		return nil, err
	}
	user, err := proxyFor(userURL, "/api/v1/users")
	if err != nil {
		return nil, err
	}
	cat, err := proxyFor(catalogURL, "/api/v1/catalog")
	if err != nil {
		return nil, err
	}
	playlist, err := proxyFor(playlistURL, "/api/v1/playlists")
	if err != nil {
		return nil, err
	}
	library, err := proxyFor(libraryURL, "/api/v1/library")
	if err != nil {
		return nil, err
	}
	playback, err := proxyFor(playbackURL, "/api/v1/playback")
	if err != nil {
		return nil, err
	}
	media, err := proxyFor(mediaURL, "/api/v1/media")
	if err != nil {
		return nil, err
	}
	events, err := proxyFor(eventsURL, "/api/v1/events")
	if err != nil {
		return nil, err
	}
	search, err := proxyFor(searchURL, "/api/v1/search")
	if err != nil {
		return nil, err
	}
	reco, err := proxyFor(recoURL, "/api/v1/recommendations")
	if err != nil {
		return nil, err
	}
	provider, err := proxyFor(providerURL, "/api/v1/providers")
	if err != nil {
		return nil, err
	}
	return &Router{
		authProxy:     auth,
		userProxy:     user,
		catalogProxy:  cat,
		playlistProxy: playlist,
		libraryProxy:  library,
		playbackProxy: playback,
		mediaProxy:    media,
		eventsProxy:   events,
		searchProxy:   search,
		recoProxy:     reco,
		providerProxy: provider,
	}, nil
}

func proxyFor(raw, basePath string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.Director
	proxy.Director = func(req *http.Request) {
		original(req)
		// preserve original host for upstream
		req.Host = target.Host

		// Ensure a single basePath prefix (avoid double-prefix bugs that cause 404).
		prefix := strings.TrimSuffix(basePath, "/")
		if !strings.HasPrefix(req.URL.Path, prefix) {
			req.URL.Path = prefix + req.URL.Path
			req.URL.RawPath = req.URL.Path
		}
	}
	return proxy, nil
}

func (r *Router) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", httpx.Health)

	api := router.Group("/api/v1")
	{
		// Explicit base + wildcard to avoid 404 when path has no trailing slash.
		api.Any("/auth", r.forward(r.authProxy))
		api.Any("/auth/*path", r.forward(r.authProxy))

		api.Any("/users", r.forward(r.userProxy))
		api.Any("/users/*path", r.forward(r.userProxy))

		api.Any("/catalog", r.forward(r.catalogProxy))
		api.Any("/catalog/*path", r.forward(r.catalogProxy))

		api.Any("/playlists", r.forward(r.playlistProxy))
		api.Any("/playlists/*path", r.forward(r.playlistProxy))

		api.Any("/library", r.forward(r.libraryProxy))
		api.Any("/library/*path", r.forward(r.libraryProxy))

		api.Any("/playback", r.forward(r.playbackProxy))
		api.Any("/playback/*path", r.forward(r.playbackProxy))

		api.Any("/media", r.forward(r.mediaProxy))
		api.Any("/media/*path", r.forward(r.mediaProxy))

		api.Any("/events", r.forward(r.eventsProxy))
		api.Any("/events/*path", r.forward(r.eventsProxy))

		api.Any("/search", r.forward(r.searchProxy))
		api.Any("/search/*path", r.forward(r.searchProxy))

		api.Any("/recommendations", r.forward(r.recoProxy))
		api.Any("/recommendations/*path", r.forward(r.recoProxy))

		api.Any("/providers", r.forward(r.providerProxy))
		api.Any("/providers/*path", r.forward(r.providerProxy))
	}
}

func (r *Router) forward(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure request ID is forwarded
		if rid, exists := c.Get(middleware.RequestIDKey); exists {
			c.Request.Header.Set("X-Request-ID", rid.(string))
		}
		// Debug path echo for quick 404 investigation (dev only; behind gateway)
		c.Writer.Header().Set("X-Gateway-Route", c.FullPath())
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
