package consumer

import (
	_ "unibee/internal/consumer/gateway"
	_ "unibee/internal/consumer/invoice"
	_ "unibee/internal/consumer/merchant"
	_ "unibee/internal/consumer/mock"
	_ "unibee/internal/consumer/payment"
	_ "unibee/internal/consumer/refund"
	_ "unibee/internal/consumer/subscription"
	_ "unibee/internal/consumer/subscription_pending_update"
	_ "unibee/internal/consumer/user"
	_ "unibee/internal/consumer/webhook"
)
