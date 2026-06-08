package service

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Trend string

const (
	GrowthTrend Trend = "Up"
	FallTrend   Trend = "Down"
	StableTrend Trend = "Stay"
)

type Activity string

const (
	FOMOActivity    Activity = "fomo"
	PanicActivity   Activity = "panic"
	ExhaustActivity Activity = "exhaust"
	WhaleActivity   Activity = "whale"
)

type Asset struct {
	ID        string
	Symbol    string
	Fullname  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	Activity
	StartedAt time.Time
	EndAt     time.Time
}

type CachedPrice struct {
	Id        string
	Symbol    string
	Price     float64
	PrevPrice float64

	BuyVolume  float64
	SellVolume float64

	Trend
	*Event
}

type AssetProvider interface {
	GetLastPrices(ctx context.Context) ([]*CachedPrice, error)
}

type PriceSaver interface {
	SaveNewPrice(ctx context.Context, assetId string, price, change, changePercent float64) error
}

type PricePublisher interface {
	PublishPrice(ctx context.Context, assetID, symbol string, price, change, changePercent float64) error
}

type priceUpdate struct {
	AssetID       string
	Price         float64
	Change        float64
	ChangePercent float64
}

type PriceGenerator struct {
	AssetProvider
	PriceSaver
	PricePublisher
	cache   map[string]*CachedPrice
	mu      sync.RWMutex
	workers int
	jobs    chan *CachedPrice
	results chan priceUpdate
}

func NewPriceGenerator(assetProvider AssetProvider, priceSaver PriceSaver, publisher PricePublisher, workers int) *PriceGenerator {
	return &PriceGenerator{
		AssetProvider:  assetProvider,
		PriceSaver:     priceSaver,
		PricePublisher: publisher,
		cache:          make(map[string]*CachedPrice),
		workers:        workers,
		jobs:           make(chan *CachedPrice, 100),
		results:        make(chan priceUpdate, 100),
	}
}

func (g *PriceGenerator) loadAssets(ctx context.Context) {
	assets, err := g.AssetProvider.GetLastPrices(ctx)
	if err != nil {
		log.Printf("generator: failed to load assets: %v", err)
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, a := range assets {
		if cached, ok := g.cache[a.Id]; ok {
			cached.Symbol = a.Symbol
		} else {
			a.PrevPrice = a.Price
			a.Trend = StableTrend
			g.cache[a.Id] = a
		}
	}

	activeIDs := make(map[string]bool)
	for _, a := range assets {
		activeIDs[a.Id] = true
	}
	for id := range g.cache {
		if !activeIDs[id] {
			delete(g.cache, id)
			log.Printf("generator: removed inactive asset from cache: %s", id)
		}
	}
}

func (g *PriceGenerator) Start(ctx context.Context) {
	log.Printf("generator: starting %d workers", g.workers)

	for i := 0; i < g.workers; i++ {
		go g.worker(ctx, i)
	}
	go g.saver(ctx)

	g.loadAssets(ctx)

	log.Printf("generator: loaded %d assets, starting ticks", len(g.cache))

	priceTicker := time.NewTicker(2 * time.Second)
	reloadTicker := time.NewTicker(5 * time.Second)
	defer priceTicker.Stop()
	defer reloadTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(g.jobs)
			close(g.results)
			log.Println("generator: stopped")
			return
		case <-priceTicker.C:
			g.dispatch(ctx)
		case <-reloadTicker.C:
			g.loadAssets(ctx)
		}
	}
}

func (g *PriceGenerator) dispatch(ctx context.Context) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, cp := range g.cache {
		select {
		case g.jobs <- cp:
		case <-ctx.Done():
			return
		}
	}
}

func (g *PriceGenerator) worker(ctx context.Context, id int) {
	for cp := range g.jobs {
		cp.PrevPrice = cp.Price
		g.updatePrice(cp)

		change := cp.Price - cp.PrevPrice
		changePercent := 0.0
		if cp.PrevPrice > 0 {
			changePercent = (change / cp.PrevPrice) * 100
		}

		select {
		case g.results <- priceUpdate{
			AssetID:       cp.Id,
			Price:         cp.Price,
			Change:        change,
			ChangePercent: changePercent,
		}:
		case <-ctx.Done():
			return
		}
	}
}

func (g *PriceGenerator) saver(ctx context.Context) {
	for result := range g.results {
		if err := g.PriceSaver.SaveNewPrice(ctx, result.AssetID, result.Price, result.Change, result.ChangePercent); err != nil {
			log.Printf("generator: failed to save price: %v", err)
		}
		g.mu.RLock()
		cp := g.cache[result.AssetID]
		g.mu.RUnlock()
		if cp != nil && g.PricePublisher != nil {
			g.PricePublisher.PublishPrice(ctx, cp.Id, cp.Symbol, result.Price, result.Change, result.ChangePercent)
		}
	}
}

func (g *PriceGenerator) updatePrice(cp *CachedPrice) {
	baseVolume := 100.0 + rand.Float64()*900.0

	switch cp.Trend {
	case GrowthTrend:
		cp.BuyVolume = baseVolume * (1.2 + rand.Float64()*0.8)
		cp.SellVolume = baseVolume * (0.3 + rand.Float64()*0.5)
	case FallTrend:
		cp.BuyVolume = baseVolume * (0.3 + rand.Float64()*0.5)
		cp.SellVolume = baseVolume * (1.2 + rand.Float64()*0.8)
	case StableTrend:
		cp.BuyVolume = baseVolume * (0.7 + rand.Float64()*0.6)
		cp.SellVolume = baseVolume * (0.7 + rand.Float64()*0.6)
	}

	if cp.Event != nil && cp.Event.EndAt.After(time.Now()) {
		switch cp.Event.Activity {
		case FOMOActivity:
			cp.BuyVolume *= 3.0
			cp.SellVolume *= 0.2
		case PanicActivity:
			cp.BuyVolume *= 0.2
			cp.SellVolume *= 3.0
		case WhaleActivity:

			if rand.Float64() > 0.5 {
				cp.BuyVolume *= 5.0
			} else {
				cp.SellVolume *= 5.0
			}
		case ExhaustActivity:
			cp.BuyVolume *= 0.5
			cp.SellVolume *= 0.5
		}
	} else {
		cp.Event = nil
	}

	volumeRatio := cp.BuyVolume / cp.SellVolume

	if volumeRatio > 1 {

		priceChange := (volumeRatio - 1) * 0.01
		if priceChange > 0.05 {
			priceChange = 0.05
		}
		cp.Price *= (1 + priceChange)
	} else if volumeRatio < 1 {
		priceChange := (1 - volumeRatio) * 0.01
		if priceChange > 0.05 {
			priceChange = 0.05
		}
		cp.Price *= (1 - priceChange)
	}

	if rand.Float64() < 0.1 {
		cp.RandomTrend()
	}

	if rand.Float64() < 0.03 {
		cp.RandomEvent()
	}

	if cp.Price < 0.0000000001 {
		cp.Price = 0.0000000001
	}
}

func (cp *CachedPrice) RandomTrend() {
	switch rand.Intn(3) {
	case 0:
		cp.Trend = GrowthTrend
	case 1:
		cp.Trend = FallTrend
	default:
		cp.Trend = StableTrend
	}
}

func (cp *CachedPrice) RandomEvent() {
	switch rand.Intn(4) {
	case 0:
		cp.Event = &Event{Activity: FOMOActivity, StartedAt: time.Now(), EndAt: time.Now().Add(time.Second * 5)}
	case 1:
		cp.Event = &Event{Activity: PanicActivity, StartedAt: time.Now(), EndAt: time.Now().Add(time.Second * 6)}
	case 2:
		cp.Event = &Event{Activity: ExhaustActivity, StartedAt: time.Now(), EndAt: time.Now().Add(time.Second * 15)}
	default:
		cp.Event = &Event{Activity: WhaleActivity, StartedAt: time.Now(), EndAt: time.Now().Add(time.Second * 3)}
	}
}

func (g *PriceGenerator) GetCache() []*CachedPrice {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]*CachedPrice, 0, len(g.cache))
	for _, cp := range g.cache {
		result = append(result, cp)
	}
	return result
}
