package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
)

// Server API 伺服器
type Server struct {
	httpServer *http.Server
	factory    *TaskFactory
}

// NewServer 建立新的 API 伺服器
func NewServer(ctx context.Context) *Server {
	cfg := config.C().API

	factory := NewTaskFactory(ctx)
	handlers := NewHandlers(factory)

	// 設定路由
	mux := http.NewServeMux()

	// 健康檢查
	mux.HandleFunc("/health", handlers.HealthCheckHandler)

	// API v1 路由
	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListTasksHandler(w, r)
		case http.MethodPost:
			handlers.CreateTaskHandler(w, r)
		default:
			MethodNotAllowedHandler(w, r)
		}
	})
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		// 根據方法和路徑分發
		switch r.Method {
		case http.MethodGet:
			handlers.GetTaskHandler(w, r)
		case http.MethodDelete:
			handlers.CancelTaskHandler(w, r)
		default:
			MethodNotAllowedHandler(w, r)
		}
	})
	mux.HandleFunc("/api/v1/storages", handlers.ListStoragesHandler)
	mux.HandleFunc("/api/v1/task-types", handlers.GetTaskTypesHandler)

	// 404 處理
	mux.HandleFunc("/", NotFoundHandler)

	// Apply middleware chain.
	var handler http.Handler = mux

	// Apply auth middleware when a token is configured.
	token := cfg.Token
	if token != "" {
		handler = AuthMiddleware()(handler)
	}

	// Add logging middleware.
	handler = loggingMiddleware(handler)

	// Add recovery middleware.
	handler = recoveryMiddleware(handler)

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      handler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		factory: factory,
	}
}

// Start 啟動伺服器
func (s *Server) Start(ctx context.Context) error {
	logger := log.FromContext(ctx).With("module", "api")

	logger.Infof("Starting API server on %s", s.httpServer.Addr)

	// 在 goroutine 中啟動伺服器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("API server error: %v", err)
		}
	}()

	// 監聽 context 取消
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("API server shutdown error: %v", err)
		}
	}()

	return nil
}

// loggingMiddleware 日誌中間件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包裝 ResponseWriter 以取得狀態碼
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Infof("%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

// recoveryMiddleware 恢復中間件
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Panic recovered: %v", err)
				WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// responseWriter 包裝 http.ResponseWriter 以擷取狀態碼
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Start initializes and starts the API server. It refuses to start without a
// token, since an open download proxy is a security risk.
func Start(ctx context.Context) error {
	cfg := config.C().API

	if !cfg.Enable {
		return nil
	}

	if cfg.Token == "" {
		return fmt.Errorf("API server is enabled but no token is set; refusing to start insecurely")
	}

	server := NewServer(ctx)
	if err := server.Start(ctx); err != nil {
		return err
	}
	StartCleanupLoop(ctx)
	return nil
}
