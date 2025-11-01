// StockExchange_LLD.go
// Low-Level Design (LLD) of a simplified stock exchange in Go.
// Targets two types of traders: Institutional and Individual.
// Institutional traders receive real-time notifications for any buy/sell orders
// of stocks they have subscribed to.
//
// This file contains a runnable in-memory prototype demonstrating:
// 1. Trader registration (institutional / individual)
// 2. Subscriptions for institutional traders
// 3. Order submission (limit orders)
// 4. Per-symbol orderbook + matching engine (single goroutine per symbol)
// 5. Notification service using Go channels for real-time delivery to institutional traders
// 6. Simple HTTP API to interact with the exchange
//
// Notes / Assumptions:
// - This is an in-memory prototype for clarity. Production needs durable storage, auth,
//   rate limiting, order lifecycle persistence, auditing, complex order types, etc.
// - Matching is price-time priority for limit orders. No market orders are implemented
//   (but easy to add).
// - Each symbol has its own goroutine that serializes all orderbook mutations (avoids heavy locks).
// - Notifications are pushed to registered institutional traders via a buffered channel.
// - For simplicity, traders are identified by string IDs.

package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

// ------------------ Domain Types ------------------

type TraderType int

const (
	Individual TraderType = iota
	Institutional
)

type Trader struct {
	ID   string     `json:"id"`
	Type TraderType `json:"type"`
	// For institutional traders we keep a notification channel.
	NotifyCh chan Notification `json:"-"`
}

// OrderSide: Buy or Sell
type Side int

const (
	Buy Side = iota
	Sell
)

func (s Side) String() string {
	if s == Buy {
		return "BUY"
	}
	return "SELL"
}

// Order represents a limit order
type Order struct {
	ID        string    `json:"id"`
	TraderID  string    `json:"trader_id"`
	Symbol    string    `json:"symbol"`
	Side      Side      `json:"side"`
	Price     float64   `json:"price"`
	Quantity  int64     `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}

// Trade represents an executed trade
type Trade struct {
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
}

// Notification for institutional traders
type Notification struct {
	Symbol string `json:"symbol"`
	Type   string `json:"type"` // "NEW_ORDER", "TRADE"
	Data   any    `json:"data"`
}

// ------------------ Registry & Notification Service ------------------

// TraderRegistry keeps track of traders and subscriptions
type TraderRegistry struct {
	mu sync.RWMutex
	// traders by id
	traders map[string]*Trader
	// subscriptions: symbol -> set of institutional trader ids
	subscriptions map[string]map[string]struct{}
}

func NewTraderRegistry() *TraderRegistry {
	return &TraderRegistry{
		traders:       make(map[string]*Trader),
		subscriptions: make(map[string]map[string]struct{}),
	}
}

func (r *TraderRegistry) RegisterTrader(id string, ttype TraderType) *Trader {
	r.mu.Lock()
	defer r.mu.Unlock()
	tr := &Trader{ID: id, Type: ttype}
	if ttype == Institutional {
		// buffered channel for notifications; size chosen arbitrarily
		tr.NotifyCh = make(chan Notification, 128)
	}
	r.traders[id] = tr
	return tr
}

func (r *TraderRegistry) GetTrader(id string) (*Trader, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tr, ok := r.traders[id]
	return tr, ok
}

// SubscribeInstitutional registers an institutional trader for symbol notifications
func (r *TraderRegistry) SubscribeInstitutional(traderID, symbol string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tr, ok := r.traders[traderID]
	if !ok {
		return fmt.Errorf("trader not found")
	}
	if tr.Type != Institutional {
		return fmt.Errorf("only institutional traders can subscribe")
	}
	set := r.subscriptions[symbol]
	if set == nil {
		set = make(map[string]struct{})
		r.subscriptions[symbol] = set
	}
	set[traderID] = struct{}{}
	return nil
}

func (r *TraderRegistry) GetSubscribers(symbol string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	set := r.subscriptions[symbol]
	if set == nil {
		return nil
	}
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

// NotifySubscribers sends notification to all institutional subscribers for a symbol
func (r *TraderRegistry) NotifySubscribers(symbol string, note Notification) {
	subs := r.GetSubscribers(symbol)
	for _, tid := range subs {
		tr, ok := r.GetTrader(tid)
		if !ok {
			continue
		}
		// non-blocking send so slow consumers don't block the engine
		select {
		case tr.NotifyCh <- note:
		default:
			// drop or metric in real system
			log.Printf("dropping notification for trader %s symbol %s", tid, symbol)
		}
	}
}

// ------------------ Orderbook & Matching Engine ------------------

// We will use two priority queues implemented using container/heap
// For buys: highest price first, then earlier timestamp
// For sells: lowest price first, then earlier timestamp

type OrderPQ []*Order

func (pq OrderPQ) Len() int      { return len(pq) }
func (pq OrderPQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

// For buys: higher price better
func (pq OrderPQ) Less(i, j int) bool {
	if pq[i].Price == pq[j].Price {
		return pq[i].Timestamp.Before(pq[j].Timestamp)
	}
	return pq[i].Price > pq[j].Price
}

func (pq *OrderPQ) Push(x any) {
	*pq = append(*pq, x.(*Order))
}

func (pq *OrderPQ) Pop() any {
	op := (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return op
}

// Helper to peek top
func (pq OrderPQ) Peek() *Order {
	if len(pq) == 0 {
		return nil
	}
	return pq[0]
}

// OrderBook for a symbol: has buy & sell PQs
type OrderBook struct {
	Symbol string
	buys   OrderPQ
	sells  OrderPQ
	mu     sync.RWMutex // only for snapshot reads; mutations are serialized by engine goroutine
}

func NewOrderBook(symbol string) *OrderBook {
	ob := &OrderBook{Symbol: symbol}
	heap.Init(&ob.buys)
	heap.Init(&ob.sells)
	return ob
}

// Snapshot for read-only views
func (ob *OrderBook) Snapshot() (buys []Order, sells []Order) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	for _, o := range ob.buys {
		buys = append(buys, *o)
	}
	for _, o := range ob.sells {
		sells = append(sells, *o)
	}
	// sort to deterministic order for consumers
	sort.Slice(buys, func(i, j int) bool {
		return buys[i].Price > buys[j].Price || (buys[i].Price == buys[j].Price && buys[i].Timestamp.Before(buys[j].Timestamp))
	})
	sort.Slice(sells, func(i, j int) bool {
		return sells[i].Price < sells[j].Price || (sells[i].Price == sells[j].Price && sells[i].Timestamp.Before(sells[j].Timestamp))
	})
	return
}

// ------------------ Symbol Engine ------------------

// Commands sent to per-symbol engine
type symbolCmd struct {
	order *Order
	resp  chan any // can be used for ack or trade list
}

// SymbolEngine serializes updates per symbol and matches orders
type SymbolEngine struct {
	symbol    string
	orderBook *OrderBook
	inbox     chan symbolCmd
	registry  *TraderRegistry
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewSymbolEngine(symbol string, registry *TraderRegistry) *SymbolEngine {
	ctx, cancel := context.WithCancel(context.Background())
	se := &SymbolEngine{
		symbol:    symbol,
		orderBook: NewOrderBook(symbol),
		inbox:     make(chan symbolCmd, 1024),
		registry:  registry,
		ctx:       ctx,
		cancel:    cancel,
	}
	go se.run()
	return se
}

func (se *SymbolEngine) Stop() { se.cancel() }

func (se *SymbolEngine) SubmitOrder(o *Order) ([]Trade, error) {
	resp := make(chan any)
	cmd := symbolCmd{order: o, resp: resp}
	select {
	case se.inbox <- cmd:
		res := <-resp
		switch v := res.(type) {
		case []Trade:
			return v, nil
		case error:
			return nil, v
		default:
			return nil, nil
		}
	case <-se.ctx.Done():
		return nil, fmt.Errorf("engine stopped")
	}
}

func (se *SymbolEngine) run() {
	for {
		select {
		case cmd := <-se.inbox:
			trades := se.matchAndBook(cmd.order)
			// notify institutional subscribers about the new order
			note := Notification{Symbol: se.symbol, Type: "NEW_ORDER", Data: cmd.order}
			se.registry.NotifySubscribers(se.symbol, note)
			// also notify trades if any
			for _, t := range trades {
				ntrade := Notification{Symbol: se.symbol, Type: "TRADE", Data: t}
				se.registry.NotifySubscribers(se.symbol, ntrade)
			}
			// respond
			cmd.resp <- trades
		case <-se.ctx.Done():
			return
		}
	}
}

// matchAndBook does simple matching; returns executed trades (possibly partial fills)
func (se *SymbolEngine) matchAndBook(o *Order) []Trade {
	ob := se.orderBook
	var trades []Trade
	// while there's remaining qty try to match
	remaining := o.Quantity
	if remaining <= 0 {
		return nil
	}

	if o.Side == Buy {
		// match with best sells (lowest price)
		for remaining > 0 && ob.sells.Len() > 0 {
			top := ob.sells.Peek()
			if top == nil {
				break
			}
			// price check: buy price must be >= sell price
			if o.Price < top.Price {
				break
			}
			// execute trade at top.Price
			tradeQty := minInt64(remaining, top.Quantity)
			trade := Trade{BuyOrderID: o.ID, SellOrderID: top.ID, Symbol: o.Symbol, Price: top.Price, Quantity: tradeQty, Timestamp: time.Now()}
			trades = append(trades, trade)
			remaining -= tradeQty
			top.Quantity -= tradeQty
			if top.Quantity == 0 {
				heap.Pop(&ob.sells)
			}
		}
		// if leftover, insert into buy book
		if remaining > 0 {
			newOrder := *o
			newOrder.Quantity = remaining
			heap.Push(&ob.buys, &newOrder)
			// update snapshot lock
			ob.mu.Lock()
			ob.mu.Unlock()
		}
	} else {
		// Sell side: match with best buys (highest price)
		for remaining > 0 && ob.buys.Len() > 0 {
			top := ob.buys.Peek()
			if top == nil {
				break
			}
			// price check: sell price must be <= buy price
			if o.Price > top.Price {
				break
			}
			tradeQty := minInt64(remaining, top.Quantity)
			trade := Trade{BuyOrderID: top.ID, SellOrderID: o.ID, Symbol: o.Symbol, Price: top.Price, Quantity: tradeQty, Timestamp: time.Now()}
			trades = append(trades, trade)
			remaining -= tradeQty
			top.Quantity -= tradeQty
			if top.Quantity == 0 {
				heap.Pop(&ob.buys)
			}
		}
		if remaining > 0 {
			newOrder := *o
			newOrder.Quantity = remaining
			heap.Push(&ob.sells, &newOrder)
			ob.mu.Lock()
			ob.mu.Unlock()
		}
	}
	return trades
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// ------------------ Exchange Coordinator ------------------

// Exchange holds many symbol engines and registry
type Exchange struct {
	registry *TraderRegistry
	engines  map[string]*SymbolEngine
	mu       sync.RWMutex
}

func NewExchange() *Exchange {
	return &Exchange{
		registry: NewTraderRegistry(),
		engines:  make(map[string]*SymbolEngine),
	}
}

func (ex *Exchange) getOrCreateEngine(symbol string) *SymbolEngine {
	ex.mu.RLock()
	eng, ok := ex.engines[symbol]
	ex.mu.RUnlock()
	if ok {
		return eng
	}
	ex.mu.Lock()
	defer ex.mu.Unlock()
	// double-check
	eng, ok = ex.engines[symbol]
	if ok {
		return eng
	}
	eng = NewSymbolEngine(symbol, ex.registry)
	ex.engines[symbol] = eng
	return eng
}

func (ex *Exchange) SubmitOrder(o *Order) ([]Trade, error) {
	eng := ex.getOrCreateEngine(o.Symbol)
	return eng.SubmitOrder(o)
}

// simple snapshot endpoint
func (ex *Exchange) Snapshot(symbol string) (buys []Order, sells []Order) {
	ex.mu.RLock()
	eng, ok := ex.engines[symbol]
	ex.mu.RUnlock()
	if !ok {
		return nil, nil
	}
	return eng.orderBook.Snapshot()
}

// ------------------ HTTP API (simple) ------------------

var ex = NewExchange()

func main() {
	// seed some traders for demo
	ex.registry.RegisterTrader("trader_ind_1", Individual)
	inst := ex.registry.RegisterTrader("trader_inst_1", Institutional)
	// spawn goroutine to read notifications for the institutional trader
	go func(tr *Trader) {
		for note := range tr.NotifyCh {
			b, _ := json.Marshal(note)
			log.Printf("[NOTIFY %s] %s\n", tr.ID, string(b))
		}
	}(inst)

	// subscribe inst to symbol ABC
	ex.registry.SubscribeInstitutional("trader_inst_1", "ABC")

	// http handlers
	http.HandleFunc("/place_order", placeOrderHandler)
	http.HandleFunc("/register_trader", registerTraderHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/snapshot", snapshotHandler)

	log.Println("Starting exchange on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// placeOrderHandler expects JSON order (without id, timestamp)
func placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TraderID string  `json:"trader_id"`
		Symbol   string  `json:"symbol"`
		Side     string  `json:"side"`
		Price    float64 `json:"price"`
		Quantity int64   `json:"quantity"`
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	side := Buy
	if req.Side == "SELL" || req.Side == "sell" {
		side = Sell
	}
	id := generateID("ORD")
	o := &Order{ID: id, TraderID: req.TraderID, Symbol: req.Symbol, Side: side, Price: req.Price, Quantity: req.Quantity, Timestamp: time.Now()}
	trades, err := ex.SubmitOrder(o)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"order_id": id, "trades": trades})
}

func registerTraderHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   string `json:"id"`
		Type string `json:"type"` // "INDIVIDUAL" or "INSTITUTIONAL"
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	typeVal := Individual
	if req.Type == "INSTITUTIONAL" {
		typeVal = Institutional
	}
	tr := ex.registry.RegisterTrader(req.ID, typeVal)
	json.NewEncoder(w).Encode(map[string]any{"trader": tr})
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TraderID string `json:"trader_id"`
		Symbol   string `json:"symbol"`
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := ex.registry.SubscribeInstitutional(req.TraderID, req.Symbol); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("ok"))
}

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	sym := r.URL.Query().Get("symbol")
	if sym == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("symbol required"))
		return
	}
	buys, sells := ex.Snapshot(sym)
	json.NewEncoder(w).Encode(map[string]any{"buys": buys, "sells": sells})
}

// ------------------ Utilities ------------------

func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), rand.Intn(1000))
}

// ------------------ End ------------------

// To run demo:
// 1. go run StockExchange_LLD.go
// 2. Register traders via POST /register_trader
// 3. Institutional traders should POST /subscribe with symbol(s) to receive notifications
// 4. Place orders via POST /place_order

// Example cURL to place buy order:
// curl -X POST localhost:8080/place_order -d '{"trader_id":"trader_ind_1","symbol":"ABC","side":"BUY","price":100.5,"quantity":10}' -H 'Content-Type: application/json'

// Limitations and next steps for production readiness:
// - Add persistence (write-ahead log, DB)
// - Authentication/Authorization
// - Order cancellation / amendment
// - Advanced order types: IOC, FOK, market orders, iceberg
// - Matching across venues, order router
// - Rate limiting + backpressure controls
// - Monitoring/metrics and exactly-once notification guarantees
// - Horizontal scaling: shard symbols across machines and use durable message queues
// - WebSocket or SSE for institutional notification delivery for real clients

// Small helper to satisfy linter
func init() {
	_ = math.MaxFloat64
	_ = strconv.Itoa
}
