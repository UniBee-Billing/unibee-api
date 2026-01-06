package statistics

import (
	"context"
	"fmt"
	"time"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"

	"github.com/gogf/gf/v2/frame/g"
)

type MerchantRealTimeCount struct {
	PlanCount               int64 `json:"planCount"`
	ActivePlanCount         int64 `json:"activePlanCount"`
	DiscountCodeCount       int64 `json:"discountCodeCount"`
	ActiveDiscountCodeCount int64 `json:"activeDiscountCodeCount"`
	GatewayCount            int64 `json:"gatewayCount"`
	InvoiceCount            int64 `json:"invoiceCount"`
	PaidInvoiceCount        int64 `json:"paidInvoiceCount"`
	RefundInvoiceCount      int64 `json:"refundInvoiceCount"`
	SubscriptionCount       int64 `json:"subscriptionCount"`
	ActiveSubscriptionCount int64 `json:"activeSubscriptionCount"`
	UserCount               int64 `json:"userCount"`
	ActiveUserCount         int64 `json:"activeUserCount"`
	LatestMonthRegisterUser int64 `json:"latestMonthRegisterUser"`
	TotalEmailSendCount     int64 `json:"totalEmailSendCount"`
	VatValidationCount      int64 `json:"vatValidationCount"`
}

type MerchantRealtimeStats struct {
	Merchant   *bean.Merchant         `json:"merchant"`
	Active     bool                   `json:"active"`
	Counts     *MerchantRealTimeCount `json:"counts"`
	LastUpdate int64                  `json:"lastUpdate"` // Last update timestamp
}

// Cache key constants
const (
	// Merchant statistics cache (15 days expiration)
	CacheKeyMerchantStats = "unibee:stats:merchant:%d:realtime"
	// Statistics update lock (prevent concurrent updates)
	CacheKeyStatsUpdateLock = "unibee:stats:merchant:%d:lock"
	// Last update time
	CacheKeyStatsLastUpdate = "unibee:stats:merchant:%d:last_update"
)

// Cache TTL constants
const (
	// Cache expiration time: 15 days
	CacheTTL = 15 * 24 * 60 * 60 // 15 days in seconds
	// Lock expiration time: 1 minute
	LockTTL = 60 // 1 minute in seconds
)

// StatsCacheService handles caching of merchant statistics
type StatsCacheService struct{}

// GetMerchantStats retrieves merchant statistics from cache or database
func (s *StatsCacheService) GetMerchantStats(ctx context.Context, merchantId uint64) (*MerchantRealtimeStats, error) {
	key := fmt.Sprintf(CacheKeyMerchantStats, merchantId)

	// Try to get from Redis cache
	cached, err := g.Redis().Get(ctx, key)
	if err == nil && !cached.IsNil() {
		var stats MerchantRealtimeStats
		if err := cached.Scan(&stats); err == nil {
			return &stats, nil
		}
	}

	// Cache miss, calculate from database and update cache
	stats := calculateMerchantStatsFromDB(ctx, merchantId)
	s.UpdateMerchantStats(ctx, merchantId, stats)

	return stats, nil
}

// UpdateMerchantStats updates merchant statistics in cache
func (s *StatsCacheService) UpdateMerchantStats(ctx context.Context, merchantId uint64, stats *MerchantRealtimeStats) error {
	key := fmt.Sprintf(CacheKeyMerchantStats, merchantId)

	// Set last update timestamp
	stats.LastUpdate = time.Now().Unix()

	// Update statistics data with 15 days expiration
	err := g.Redis().SetEX(ctx, key, stats, CacheTTL)
	if err != nil {
		return err
	}

	// Update last update time with 15 days expiration (for backward compatibility)
	lastUpdateKey := fmt.Sprintf(CacheKeyStatsLastUpdate, merchantId)
	g.Redis().SetEX(ctx, lastUpdateKey, stats.LastUpdate, CacheTTL)

	return nil
}

// IncrementCount atomically increments a specific count field
func (s *StatsCacheService) IncrementCount(ctx context.Context, merchantId uint64, field string, delta int64) {
	key := fmt.Sprintf(CacheKeyMerchantStats, merchantId)
	lastUpdateKey := fmt.Sprintf(CacheKeyStatsLastUpdate, merchantId)
	currentTime := time.Now().Unix()

	// Use Redis HINCRBY for atomic increment
	g.Redis().HIncrBy(ctx, key, field, delta)

	// Update last update timestamp in the stats structure
	g.Redis().HSet(ctx, key, map[string]interface{}{
		"lastUpdate": currentTime,
	})

	// Reset expiration time to 15 days
	g.Redis().Expire(ctx, key, CacheTTL)

	// Update last update time with 15 days expiration (for backward compatibility)
	g.Redis().SetEX(ctx, lastUpdateKey, currentTime, CacheTTL)
}

// UpdateMerchantStatsAsync updates merchant statistics asynchronously with lock
func (s *StatsCacheService) UpdateMerchantStatsAsync(ctx context.Context, merchantId uint64) {
	lockKey := fmt.Sprintf(CacheKeyStatsUpdateLock, merchantId)

	// Try to acquire lock with 1 minute expiration
	locked, _ := g.Redis().SetNX(ctx, lockKey, "1")
	if locked {
		// Set lock expiration to 1 minute
		g.Redis().Expire(ctx, lockKey, LockTTL)
		defer g.Redis().Del(ctx, lockKey)

		// Calculate statistics from database
		stats := calculateMerchantStatsFromDB(ctx, merchantId)

		// Update cache (this will reset expiration to 15 days)
		s.UpdateMerchantStats(ctx, merchantId, stats)
	}
}

// GetMerchantStatsWithFallback gets merchant stats with fallback to database
func (s *StatsCacheService) GetMerchantStatsWithFallback(ctx context.Context, merchantId uint64) *MerchantRealtimeStats {
	// Try to get from cache first
	if stats, err := s.GetMerchantStats(ctx, merchantId); err == nil && stats != nil {
		return stats
	}

	// Cache failed, fallback to database
	return calculateMerchantStatsFromDB(ctx, merchantId)
}

// GetMerchantRealtimeStats gets merchant statistics with cache support
func GetMerchantRealtimeStats(ctx context.Context, merchantId uint64) *MerchantRealtimeStats {
	if merchantId <= 0 {
		return &MerchantRealtimeStats{}
	}

	cacheService := &StatsCacheService{}
	return cacheService.GetMerchantStatsWithFallback(ctx, merchantId)
}

// calculateMerchantStatsFromDB calculates merchant statistics directly from database
func calculateMerchantStatsFromDB(ctx context.Context, merchantId uint64) *MerchantRealtimeStats {
	if merchantId <= 0 {
		return &MerchantRealtimeStats{}
	}

	// Get merchant basic info
	merchant := query.GetMerchantById(ctx, merchantId)
	if merchant == nil {
		return &MerchantRealtimeStats{}
	}

	// Get counts using DAO queries
	counts := &MerchantRealTimeCount{}

	// Plan counts
	planCount, _ := dao.Plan.Ctx(ctx).
		Where(dao.Plan.Columns().MerchantId, merchantId).
		Where(dao.Plan.Columns().IsDeleted, 0).
		Count()
	counts.PlanCount = int64(planCount)

	activePlanCount, _ := dao.Plan.Ctx(ctx).
		Where(dao.Plan.Columns().MerchantId, merchantId).
		Where(dao.Plan.Columns().IsDeleted, 0).
		Where(dao.Plan.Columns().Status, consts.PlanStatusActive).
		Count()
	counts.ActivePlanCount = int64(activePlanCount)

	// Discount code counts
	discountCount, _ := dao.MerchantDiscountCode.Ctx(ctx).
		Where(dao.MerchantDiscountCode.Columns().MerchantId, merchantId).
		Where(dao.MerchantDiscountCode.Columns().IsDeleted, 0).
		Count()
	counts.DiscountCodeCount = int64(discountCount)

	activeDiscountCount, _ := dao.MerchantDiscountCode.Ctx(ctx).
		Where(dao.MerchantDiscountCode.Columns().MerchantId, merchantId).
		Where(dao.MerchantDiscountCode.Columns().IsDeleted, 0).
		Where(dao.MerchantDiscountCode.Columns().Status, consts.DiscountStatusActive).
		Count()
	counts.ActiveDiscountCodeCount = int64(activeDiscountCount)

	// Gateway count
	gatewayCount, _ := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
		Where(dao.MerchantGateway.Columns().IsDeleted, 0).
		Count()
	counts.GatewayCount = int64(gatewayCount)

	// Invoice counts
	invoiceCount, _ := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().MerchantId, merchantId).
		Count()
	counts.InvoiceCount = int64(invoiceCount)

	paidInvoiceCount, _ := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().MerchantId, merchantId).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusPaid).
		Count()
	counts.PaidInvoiceCount = int64(paidInvoiceCount)

	refundInvoiceCount, _ := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().MerchantId, merchantId).
		Where("refund_id IS NOT NULL AND refund_id != ''").
		Count()
	counts.RefundInvoiceCount = int64(refundInvoiceCount)

	// Subscription counts
	subscriptionCount, _ := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		Count()
	counts.SubscriptionCount = int64(subscriptionCount)

	activeSubscriptionCount, _ := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		Where(dao.Subscription.Columns().Status, consts.SubStatusActive).
		Count()
	counts.ActiveSubscriptionCount = int64(activeSubscriptionCount)

	// User counts
	userCount, _ := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Where(dao.UserAccount.Columns().IsDeleted, 0).
		Count()
	counts.UserCount = int64(userCount)

	activeUserCount, _ := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Where(dao.UserAccount.Columns().IsDeleted, 0).
		Where(dao.UserAccount.Columns().Status, 0). // 0 = Active
		Count()
	counts.ActiveUserCount = int64(activeUserCount)

	// Latest month register user count
	oneMonthAgo := time.Now().AddDate(0, -1, 0).Unix()
	latestMonthRegisterUser, _ := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Where(dao.UserAccount.Columns().IsDeleted, 0).
		WhereGTE(dao.UserAccount.Columns().CreateTime, oneMonthAgo).
		Count()
	counts.LatestMonthRegisterUser = int64(latestMonthRegisterUser)

	// Email send count
	emailSendCount, _ := dao.MerchantEmailHistory.Ctx(ctx).
		Where(dao.MerchantEmailHistory.Columns().MerchantId, merchantId).
		Count()
	counts.TotalEmailSendCount = int64(emailSendCount)

	// VAT validation count
	vatValidationCount, _ := dao.MerchantVatNumberVerifyHistory.Ctx(ctx).
		Where(dao.MerchantVatNumberVerifyHistory.Columns().MerchantId, merchantId).
		Where(dao.MerchantVatNumberVerifyHistory.Columns().IsDeleted, 0).
		Count()
	counts.VatValidationCount = int64(vatValidationCount)

	// Check if merchant is active (not deleted)
	isActive := merchant.IsDeleted == 0

	return &MerchantRealtimeStats{
		Merchant:   bean.SimplifyMerchant(merchant),
		Active:     isActive,
		Counts:     counts,
		LastUpdate: time.Now().Unix(),
	}
}

// Event-driven increment functions for real-time updates

// OnUserRegistered handles user registration event
func OnUserRegistered(ctx context.Context, user *entity.UserAccount) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, user.MerchantId, "userCount", 1)
	cacheService.IncrementCount(ctx, user.MerchantId, "latestMonthRegisterUser", 1)
}

// OnPlanStatusChanged handles plan status change event
func OnPlanStatusChanged(ctx context.Context, plan *entity.Plan, oldStatus int) {
	cacheService := &StatsCacheService{}

	if plan.Status == consts.PlanStatusActive && oldStatus != consts.PlanStatusActive {
		cacheService.IncrementCount(ctx, plan.MerchantId, "activePlanCount", 1)
	} else if plan.Status != consts.PlanStatusActive && oldStatus == consts.PlanStatusActive {
		cacheService.IncrementCount(ctx, plan.MerchantId, "activePlanCount", -1)
	}
}

// OnPlanCreated handles plan creation event
func OnPlanCreated(ctx context.Context, plan *entity.Plan) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, plan.MerchantId, "planCount", 1)
	if plan.Status == consts.PlanStatusActive {
		cacheService.IncrementCount(ctx, plan.MerchantId, "activePlanCount", 1)
	}
}

// OnPlanDeleted handles plan deletion event
func OnPlanDeleted(ctx context.Context, merchantId uint64, planId uint64) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, merchantId, "planCount", -1)
	// Note: We can't easily determine if it was active, so we'll let the periodic update handle this
}

// OnSubscriptionStatusChanged handles subscription status change event
func OnSubscriptionStatusChanged(ctx context.Context, sub *entity.Subscription, oldStatus int) {
	cacheService := &StatsCacheService{}

	if sub.Status == consts.SubStatusActive && oldStatus != consts.SubStatusActive {
		cacheService.IncrementCount(ctx, sub.MerchantId, "activeSubscriptionCount", 1)
	} else if sub.Status != consts.SubStatusActive && oldStatus == consts.SubStatusActive {
		cacheService.IncrementCount(ctx, sub.MerchantId, "activeSubscriptionCount", -1)
	}
}

// OnSubscriptionCreated handles subscription creation event
func OnSubscriptionCreated(ctx context.Context, sub *entity.Subscription) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, sub.MerchantId, "subscriptionCount", 1)
	if sub.Status == consts.SubStatusActive {
		cacheService.IncrementCount(ctx, sub.MerchantId, "activeSubscriptionCount", 1)
	}
}

// OnInvoiceStatusChanged handles invoice status change event
func OnInvoiceStatusChanged(ctx context.Context, invoice *entity.Invoice, oldStatus int) {
	cacheService := &StatsCacheService{}

	if invoice.Status == consts.InvoiceStatusPaid && oldStatus != consts.InvoiceStatusPaid {
		cacheService.IncrementCount(ctx, invoice.MerchantId, "paidInvoiceCount", 1)
	} else if invoice.Status != consts.InvoiceStatusPaid && oldStatus == consts.InvoiceStatusPaid {
		cacheService.IncrementCount(ctx, invoice.MerchantId, "paidInvoiceCount", -1)
	}
}

// OnInvoiceCreated handles invoice creation event
func OnInvoiceCreated(ctx context.Context, invoice *entity.Invoice) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, invoice.MerchantId, "invoiceCount", 1)
	if invoice.Status == consts.InvoiceStatusPaid {
		cacheService.IncrementCount(ctx, invoice.MerchantId, "paidInvoiceCount", 1)
	}
}

// OnRefundCreated handles refund creation event
func OnRefundCreated(ctx context.Context, merchantId uint64) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, merchantId, "refundInvoiceCount", 1)
}

// OnEmailSent handles email sent event
func OnEmailSent(ctx context.Context, merchantId uint64) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, merchantId, "totalEmailSendCount", 1)
}

// OnVatValidated handles VAT validation event
func OnVatValidated(ctx context.Context, merchantId uint64) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, merchantId, "vatValidationCount", 1)
}

// OnDiscountCodeStatusChanged handles discount code status change event
func OnDiscountCodeStatusChanged(ctx context.Context, discount *entity.MerchantDiscountCode, oldStatus int) {
	cacheService := &StatsCacheService{}

	if discount.Status == consts.DiscountStatusActive && oldStatus != consts.DiscountStatusActive {
		cacheService.IncrementCount(ctx, discount.MerchantId, "activeDiscountCodeCount", 1)
	} else if discount.Status != consts.DiscountStatusActive && oldStatus == consts.DiscountStatusActive {
		cacheService.IncrementCount(ctx, discount.MerchantId, "activeDiscountCodeCount", -1)
	}
}

// OnDiscountCodeCreated handles discount code creation event
func OnDiscountCodeCreated(ctx context.Context, discount *entity.MerchantDiscountCode) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, discount.MerchantId, "discountCodeCount", 1)
	if discount.Status == consts.DiscountStatusActive {
		cacheService.IncrementCount(ctx, discount.MerchantId, "activeDiscountCodeCount", 1)
	}
}

// OnGatewayCreated handles gateway creation event
func OnGatewayCreated(ctx context.Context, gateway *entity.MerchantGateway) {
	cacheService := &StatsCacheService{}
	cacheService.IncrementCount(ctx, gateway.MerchantId, "gatewayCount", 1)
}

// UpdateMerchantStatsCron updates merchant statistics for all active merchants
func UpdateMerchantStatsCron(ctx context.Context) {
	// Get all active merchants
	merchants := query.GetActiveMerchantList(ctx)

	g.Log().Infof(ctx, "UpdateMerchantStatsCron start, processing %d merchants", len(merchants))

	for i, merchant := range merchants {
		// Check if update is needed (avoid frequent updates)
		lastUpdateKey := fmt.Sprintf(CacheKeyStatsLastUpdate, merchant.Id)
		lastUpdate, _ := g.Redis().Get(ctx, lastUpdateKey)

		// Update if more than 1 hour since last update
		if lastUpdate.IsNil() || time.Now().Unix()-lastUpdate.Int64() > 3600 {
			cacheService := &StatsCacheService{}
			cacheService.UpdateMerchantStatsAsync(ctx, merchant.Id)
			g.Log().Infof(ctx, "UpdateMerchantStatsCron processed merchant %d/%d (ID: %d)", i+1, len(merchants), merchant.Id)
		} else {
			g.Log().Debugf(ctx, "UpdateMerchantStatsCron skipped merchant %d/%d (ID: %d) - recently updated", i+1, len(merchants), merchant.Id)
		}

		// Rest for 10 seconds after each merchant (except the last one)
		if i < len(merchants)-1 {
			g.Log().Debugf(ctx, "UpdateMerchantStatsCron resting 10 seconds before next merchant")
			time.Sleep(10 * time.Second)
		}
	}

	g.Log().Infof(ctx, "UpdateMerchantStatsCron completed, processed %d merchants", len(merchants))
}
