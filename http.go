package campusstore

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func NewHandler(service *Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/styles.css", assetHandler("styles.css", "text/css; charset=utf-8"))
	mux.HandleFunc("/app.js", assetHandler("app.js", "text/javascript; charset=utf-8"))
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"products": service.Products()})
	})
	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var input credentialsRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, ErrInvalidInput)
			return
		}
		member, err := service.Register(input.Username, input.Password)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"member": member})
	})
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var input credentialsRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, ErrInvalidInput)
			return
		}
		member, token, err := service.Login(input.Username, input.Password)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"member": member, "token": token})
	})
	mux.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		memberID := memberIDFromRequest(r)
		view, err := service.Cart(memberID)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, view)
	})
	mux.HandleFunc("/api/cart/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var input addItemRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, ErrInvalidInput)
			return
		}
		view, err := service.AddToCart(input.MemberID, input.ProductID, input.Quantity)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, view)
	})
	return mux
}

type credentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type addItemRequest struct {
	MemberID  string `json:"member_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	assetHandler("index.php", "text/html; charset=utf-8")(w, r)
}

func assetHandler(name, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		content, err := siteFiles.ReadFile(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	return nil
}

func memberIDFromRequest(r *http.Request) string {
	if value := strings.TrimSpace(r.URL.Query().Get("member_id")); value != "" {
		return value
	}
	return strings.TrimSpace(r.Header.Get("X-Member-ID"))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
	writeError(w, http.StatusMethodNotAllowed, errors.New("请求方法不支持"))
}

func writeServiceError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	switch {
	case errors.Is(err, ErrMemberNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrProductNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrInvalidCredentials):
		status = http.StatusUnauthorized
	case errors.Is(err, ErrUsernameTaken):
		status = http.StatusConflict
	}
	writeError(w, status, err)
}
