package setup

import (
	"context"
	"unibee/internal/logic/analysis/segment"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/operation_log"
)

func MerchantSegmentSetup(ctx context.Context, merchantId uint64, serverSideSecret string, userPortalSecret string) error {
	err := merchant_config.SetMerchantConfig(ctx, merchantId, segment.KeyMerchantSegmentServer, serverSideSecret)
	err = merchant_config.SetMerchantConfig(ctx, merchantId, segment.KeyMerchantSegmentUserPortal, userPortalSecret)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     merchantId,
		Target:         "Segment",
		Content:        "SetupSegment",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return err
	}
	return err
}

func CleanMerchantSegment(ctx context.Context, merchantId uint64) error {
	err := merchant_config.SetMerchantConfig(ctx, merchantId, segment.KeyMerchantSegmentServer, "")
	if err != nil {
		return err
	}
	return merchant_config.SetMerchantConfig(ctx, merchantId, segment.KeyMerchantSegmentUserPortal, "")
}
