package campusstore

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码不正确")
	ErrUsernameTaken      = errors.New("用户名已存在")
	ErrInvalidInput       = errors.New("请求内容不完整")
	ErrProductNotFound    = errors.New("商品不存在")
	ErrMemberNotFound     = errors.New("会员不存在")
	ErrInvalidQuantity    = errors.New("商品数量必须为正数")
)

type Product struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Image       string          `json:"image"`
}

type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Subtotal string  `json:"subtotal"`
}

type CartView struct {
	MemberID string     `json:"member_id"`
	Items    []CartItem `json:"items"`
	Total    string     `json:"total"`
}

type cart struct {
	MemberID string
	Items    map[string]int
}

type memberRecord struct {
	Member
	Password string
	Cart     *cart
}

type Service struct {
	mu       sync.RWMutex
	products []Product
	byID     map[string]Product
	members  map[string]*memberRecord
	byName   map[string]string
	sessions map[string]string
}

func NewFixture() *Service {
	products := []Product{
		{ID: "canvas-bag", Name: "校园帆布袋", Description: "轻便耐用的校园日常帆布袋", Price: decimal.NewFromInt(39), Image: "canvas-bag.svg"},
		{ID: "campus-badge", Name: "校园徽章", Description: "别在书包上的校园纪念徽章", Price: decimal.NewFromInt(12), Image: "campus-badge.svg"},
		{ID: "postcard", Name: "校园明信片", Description: "记录校园风景的明信片套装", Price: decimal.NewFromInt(16), Image: "postcard.svg"},
		{ID: "notebook", Name: "校园笔记本", Description: "适合课程记录的线圈笔记本", Price: decimal.NewFromInt(22), Image: "notebook.svg"},
	}
	byID := make(map[string]Product, len(products))
	for _, product := range products {
		byID[product.ID] = product
	}
	return &Service{
		products: products,
		byID:     byID,
		members:  make(map[string]*memberRecord),
		byName:   make(map[string]string),
		sessions: make(map[string]string),
	}
}

func (s *Service) Products() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Product(nil), s.products...)
}

func (s *Service) Register(username, password string) (Member, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return Member{}, ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byName[username]; exists {
		return Member{}, ErrUsernameTaken
	}
	id := fmt.Sprintf("member-%03d", len(s.members)+1)
	record := &memberRecord{
		Member:   Member{ID: id, Username: username},
		Password: password,
		Cart:     &cart{MemberID: id},
	}
	s.members[id] = record
	s.byName[username] = id
	return record.Member, nil
}

func (s *Service) Login(username, password string) (Member, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, exists := s.byName[strings.TrimSpace(username)]
	if !exists || s.members[id].Password != password {
		return Member{}, "", ErrInvalidCredentials
	}
	token := fmt.Sprintf("session-%03d", len(s.sessions)+1)
	s.sessions[token] = id
	return s.members[id].Member, token, nil
}

func (s *Service) MemberBySession(token string) (Member, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, exists := s.sessions[token]
	if !exists {
		return Member{}, ErrMemberNotFound
	}
	return s.members[id].Member, nil
}

func (s *Service) AddToCart(memberID, productID string, quantity int) (CartView, error) {
	if quantity <= 0 {
		return CartView{}, ErrInvalidQuantity
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record, exists := s.members[memberID]
	if !exists {
		return CartView{}, ErrMemberNotFound
	}
	if _, exists := s.byID[productID]; !exists {
		return CartView{}, ErrProductNotFound
	}
	record.Cart.Items[productID] += quantity
	return s.cartViewLocked(record.Cart), nil
}

func (s *Service) Cart(memberID string) (CartView, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, exists := s.members[memberID]
	if !exists {
		return CartView{}, ErrMemberNotFound
	}
	return s.cartViewLocked(record.Cart), nil
}

func (s *Service) cartViewLocked(value *cart) CartView {
	view := CartView{MemberID: value.MemberID, Items: make([]CartItem, 0), Total: "0.00"}
	for _, product := range s.products {
		quantity := value.Items[product.ID]
		if quantity == 0 {
			continue
		}
		subtotal := product.Price.Mul(decimal.NewFromInt(int64(quantity)))
		view.Items = append(view.Items, CartItem{Product: product, Quantity: quantity, Subtotal: subtotal.StringFixed(2)})
		view.Total = decimal.RequireFromString(view.Total).Add(subtotal).StringFixed(2)
	}
	return view
}
