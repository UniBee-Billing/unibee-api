package onetime

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/subscription_onetimeaddon"
	discount2 "unibee/internal/logic/invoice/discount"
	"unibee/internal/logic/subscription/handler"
	entity "unibee/internal/model/entity/default"
)

type Onetime struct {
}

func (i Onetime) PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Onetime) PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Onetime) PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	if payment.TotalAmount <= payment.RefundAmount {
		err := discount2.InvoiceRollbackAllDiscountsFromPayment(ctx, payment.InvoiceId, payment.PaymentId)
		if err != nil {
			fmt.Printf("InvoiceRollbackAllDiscountsFromPayment error:%s", err.Error())
		}
	}
}

func (i Onetime) PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Onetime) PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Onetime) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {

}

func (i Onetime) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	var metadata = make(map[string]interface{})
	if len(payment.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
		if err != nil {
			g.Log().Errorf(ctx, "PaymentCancelCallback Unmarshal Metadata error:%s", err.Error())
		}
	}
	if id, ok := metadata["SubscriptionOnetimeAddonId"]; ok {
		idInt, err := strconv.Atoi(fmt.Sprintf("%v", id))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentCancelCallback panic int: %s err:%s", id, err)
			return
		}
		_, err = handler.HandleOnetimeAddonPaymentCancel(ctx, uint64(idInt))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentCancelCallback HandleOnetimeAddonPaymentCancel int: %s err:%s", id, err)
			return
		}
		one := handler.SubscriptionOnetimeAddonDetail(ctx, uint64(idInt))
		if one != nil {
			subscription_onetimeaddon.SendMerchantSubscriptionOnetimeAddonWebhookBackground(payment.MerchantId, one, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ONETIME_ADDON_CANCELLED)
		}
	}
	//err := discount.UserDiscountRollbackFromPayment(ctx, payment.PaymentId)
	//if err != nil {
	//	fmt.Printf("UserDiscountRollbackFromPayment error:%s", err.Error())
	//}
}

func (i Onetime) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	var metadata = make(map[string]interface{})
	if len(payment.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
		if err != nil {
			fmt.Printf("PaymentCreateCallback Unmarshal Metadata error:%s", err.Error())
		}
	}
	if id, ok := metadata["SubscriptionOnetimeAddonId"]; ok {
		idInt, err := strconv.Atoi(fmt.Sprintf("%v", id))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentCancelCallback panic int: %s err:%s", id, err)
			return
		}
		err = handler.UpdateSubscriptionAddonPurchasePaymentId(ctx, uint64(idInt), payment.PaymentId)
		if err != nil {
			g.Log().Errorf(ctx, "PaymentCancelCallback UpdateSubscriptionAddonPurchasePaymentId int: %s err:%s", id, err)
			return
		}
		one := handler.SubscriptionOnetimeAddonDetail(ctx, uint64(idInt))
		if one != nil {
			subscription_onetimeaddon.SendMerchantSubscriptionOnetimeAddonWebhookBackground(payment.MerchantId, one, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ONETIME_ADDON_CREATED)
		}
	}
}

func (i Onetime) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	var metadata = make(map[string]interface{})
	if len(payment.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
		if err != nil {
			fmt.Printf("PaymentSuccessCallback Unmarshal Metadata error:%s", err.Error())
		}
	}
	if id, ok := metadata["SubscriptionOnetimeAddonId"]; ok {
		idInt, err := strconv.Atoi(fmt.Sprintf("%v", id))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentSuccessCallback panic int: %s err:%s", id, err)
			return
		}
		err = handler.UpdateSubscriptionAddonPurchasePaymentId(ctx, uint64(idInt), payment.PaymentId)
		if err != nil {
			g.Log().Errorf(ctx, "PaymentCancelCallback UpdateSubscriptionAddonPurchasePaymentId int: %s err:%s", id, err)
			return
		}
		_, err = handler.HandleOnetimeAddonPaymentSuccess(ctx, uint64(idInt))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentSuccessCallback HandleOnetimeAddonPaymentCancel int: %s err:%s", id, err)
			return
		}
		one := handler.SubscriptionOnetimeAddonDetail(ctx, uint64(idInt))
		if one != nil {
			subscription_onetimeaddon.SendMerchantSubscriptionOnetimeAddonWebhookBackground(payment.MerchantId, one, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ONETIME_ADDON_SUCCESS)
		}
	}
}

func (i Onetime) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	var metadata = make(map[string]interface{})
	if len(payment.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
		if err != nil {
			fmt.Printf("PaymentFailureCallback PaymentFailureCallback Unmarshal Metadata error:%s", err.Error())
		}
	}
	if id, ok := metadata["SubscriptionOnetimeAddonId"]; ok {
		idInt, err := strconv.Atoi(fmt.Sprintf("%v", id))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentFailureCallback panic int: %s err:%s", id, err)
			return
		}
		_, err = handler.HandleOnetimeAddonPaymentFailure(ctx, uint64(idInt))
		if err != nil {
			g.Log().Errorf(ctx, "PaymentFailureCallback HandleOnetimeAddonPaymentFailure int: %s err:%s", id, err)
			return
		}
		one := handler.SubscriptionOnetimeAddonDetail(ctx, uint64(idInt))
		if one != nil {
			subscription_onetimeaddon.SendMerchantSubscriptionOnetimeAddonWebhookBackground(payment.MerchantId, one, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ONETIME_ADDON_EXPIRED)
		}
	}
	//err := discount.UserDiscountRollbackFromPayment(ctx, payment.PaymentId)
	//if err != nil {
	//	fmt.Printf("UserDiscountRollbackFromPayment error:%s", err.Error())
	//}
}
