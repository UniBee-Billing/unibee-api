package api

import (
	"context"
	"testing"
	_interface "unibee/internal/interface"
)

func TestForCryptadium(t *testing.T) {
	pay := &Cryptadium{}
	// need account active
	_, _, _ = pay.GatewayTest(context.Background(), &_interface.GatewayTestReq{
		Key:                 "30zuWWwMEC",
		Secret:              "afuLxr90Gck0ksqszZCqf8ynPvx7myuTYbmyUy0WmLyG88T7xYXSQHXcuG6B3CNTuSbJE2XBjrDVKrpRb0KNabDQDJEM1zdh7eYsgJ25BVujECSCdrkC2A3SPLF5xKLwLeYrS7uT7PMEudKA1BuKPbr6nFcVD7H1K7Xjqnwf6REnzyjvAEDwbAgvzbDLeYPAGV4wFmWPQMP7WCGzgtU93XYKapwMjkAhCEfc4ypZt0eQxjXwUPpyDQHu7v",
		SubGateway:          "",
		GatewayPaymentTypes: nil,
	})

}
