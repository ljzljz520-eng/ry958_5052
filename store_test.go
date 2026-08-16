package campusstore_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	campusstore "campus-creative-store"
)

func TestCatalogShowsFourCampusProducts(t *testing.T) {
	service := campusstore.NewFixture()
	products := service.Products()
	if len(products) != 4 {
		t.Fatalf("expected four products, got %d", len(products))
	}
	for index, name := range []string{"校园帆布袋", "校园徽章", "校园明信片", "校园笔记本"} {
		if products[index].Name != name {
			t.Fatalf("expected product %q at position %d, got %q", name, index, products[index].Name)
		}
	}
}

func TestMemberCanRegisterAndLoginWithUsername(t *testing.T) {
	service := campusstore.NewFixture()
	registered, err := service.Register("林同学", "campus-pass")
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	loggedIn, token, err := service.Login("林同学", "campus-pass")
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if token == "" || loggedIn.ID != registered.ID || loggedIn.Username != "林同学" {
		t.Fatalf("login did not return the registered member")
	}
	fromSession, err := service.MemberBySession(token)
	if err != nil || fromSession.Username != "林同学" {
		t.Fatalf("session did not identify the member")
	}
}

func TestProductsEndpointReturnsCatalog(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	campusstore.NewHandler(campusstore.NewFixture()).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected products endpoint status 200, got %d", recorder.Code)
	}
	var body struct {
		Products []campusstore.Product `json:"products"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("products endpoint returned invalid JSON: %v", err)
	}
	if len(body.Products) != 4 {
		t.Fatalf("expected four products from endpoint, got %d", len(body.Products))
	}
}

func TestHomepageServesCampusStoreAndStylesheet(t *testing.T) {
	handler := campusstore.NewHandler(campusstore.NewFixture())
	home := httptest.NewRecorder()
	handler.ServeHTTP(home, httptest.NewRequest(http.MethodGet, "/", nil))
	if home.Code != http.StatusOK || !strings.Contains(home.Body.String(), "校园文创商品站") || !strings.Contains(home.Body.String(), "/styles.css") {
		t.Fatalf("homepage did not serve the campus store with its stylesheet")
	}
	stylesheet := httptest.NewRecorder()
	handler.ServeHTTP(stylesheet, httptest.NewRequest(http.MethodGet, "/styles.css", nil))
	if stylesheet.Code != http.StatusOK || !strings.Contains(stylesheet.Header().Get("Content-Type"), "text/css") {
		t.Fatalf("stylesheet was not served as CSS")
	}
}

func TestNewMemberStartsWithEmptyCart(t *testing.T) {
	service := campusstore.NewFixture()
	member, err := service.Register("吴同学", "campus-pass")
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	cart, err := service.Cart(member.ID)
	if err != nil {
		t.Fatalf("cart query returned error: %v", err)
	}
	if len(cart.Items) != 0 || cart.Total != "0.00" {
		t.Fatalf("new member cart was not empty")
	}
}

func TestNewMemberCanAddBadgeToCart(t *testing.T) {
	service := campusstore.NewFixture()
	member, err := service.Register("周同学", "campus-pass")
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("adding badge crashed: %v", recovered)
		}
	}()
	cart, err := service.AddToCart(member.ID, "campus-badge", 1)
	if err != nil {
		t.Fatalf("adding badge returned error: %v", err)
	}
	if len(cart.Items) != 1 || cart.Items[0].Product.ID != "campus-badge" || cart.Items[0].Quantity != 1 || cart.Total != "12.00" {
		t.Fatalf("badge was not added to the member cart")
	}
}

func TestCartRejectsInvalidRequests(t *testing.T) {
	service := campusstore.NewFixture()
	member, err := service.Register("赵同学", "campus-pass")
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if _, err := service.AddToCart(member.ID, "missing-product", 1); err != campusstore.ErrProductNotFound {
		t.Fatalf("expected unknown product to be rejected, got %v", err)
	}
	if _, err := service.AddToCart(member.ID, "notebook", 0); err != campusstore.ErrInvalidQuantity {
		t.Fatalf("expected zero quantity to be rejected, got %v", err)
	}
}
