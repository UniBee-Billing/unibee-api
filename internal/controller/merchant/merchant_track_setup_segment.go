package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/analysis/segment/setup"

	"unibee/api/merchant/track"
)

func (c *ControllerTrack) SetupSegment(ctx context.Context, req *track.SetupSegmentReq) (res *track.SetupSegmentRes, err error) {
	err = setup.MerchantSegmentSetup(ctx, _interface.GetMerchantId(ctx), req.ServerSideSecret, req.UserPortalSecret)
	if err != nil {
		return nil, err
	}
	return &track.SetupSegmentRes{}, nil
}
