package preload

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
)

// UserPreloadData holds all preloaded data for users
type UserPreloadData struct {
	Gateways           map[string]*entity.MerchantGateway
	UserTaxPercentages map[uint64]int64
}

func UserListPreloadForContext(ctx context.Context, users []*entity.UserAccount) {
	if len(users) > 3 {
		preloadData := UserListPreload(ctx, users)
		if preloadData != nil && _interface.GetBulkPreloadData(ctx) != nil {
			var gateways = make(map[uint64]*entity.MerchantGateway)
			for _, one := range preloadData.Gateways {
				gateways[one.Id] = one
			}
			_interface.Context().Get(ctx).PreloadData.Gateways = gateways
		}
	}
}

// UserListPreload preloads all related data for users in bulk
func UserListPreload(ctx context.Context, users []*entity.UserAccount) *UserPreloadData {
	preload := &UserPreloadData{
		Gateways:           make(map[string]*entity.MerchantGateway),
		UserTaxPercentages: make(map[uint64]int64),
	}

	if len(users) == 0 {
		return preload
	}

	// Collect unique gateway IDs
	gatewayIds := make(map[string]bool)
	userCountryPairs := make(map[string]bool)

	for _, user := range users {
		if len(user.GatewayId) > 0 {
			gatewayIds[user.GatewayId] = true
		}
		if len(user.CountryCode) > 0 {
			userCountryPairs[fmt.Sprintf("%d_%s", user.MerchantId, user.CountryCode)] = true
		}
	}

	// Bulk query gateways
	if len(gatewayIds) > 0 {
		// Convert string IDs to uint64 for database query
		var gatewayIdUint64s []uint64
		for idStr := range gatewayIds {
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				gatewayIdUint64s = append(gatewayIdUint64s, id)
			}
		}

		if len(gatewayIdUint64s) > 0 {
			var gateways []*entity.MerchantGateway
			err := dao.MerchantGateway.Ctx(ctx).WhereIn(dao.MerchantGateway.Columns().Id, gatewayIdUint64s).Scan(&gateways)
			if err == nil {
				for _, gateway := range gateways {
					// Use string ID as key to match UserAccount.GatewayId
					preload.Gateways[strconv.FormatUint(gateway.Id, 10)] = gateway
				}
			}
		}
	}

	// Bulk query VAT rates for country pairs
	if len(userCountryPairs) > 0 {
		for pair := range userCountryPairs {
			parts := strings.Split(pair, "_")
			if len(parts) == 2 {
				merchantId, _ := strconv.ParseUint(parts[0], 10, 64)
				countryCode := parts[1]

				vatRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, merchantId, countryCode)
				if err == nil && vatRate != nil {
					// Find all users with this merchant and country
					for _, user := range users {
						if user.MerchantId == merchantId && user.CountryCode == countryCode {
							preload.UserTaxPercentages[user.Id] = vatRate.StandardTaxPercentage
						}
					}
				}
			}
		}
	}

	return preload
}
