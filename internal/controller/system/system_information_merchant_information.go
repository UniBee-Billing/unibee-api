package system

import (
	"context"
	"os"
	"strings"
	"unibee-api/api/system/information"
)

func (c *ControllerInformation) MerchantInformation(ctx context.Context, req *information.MerchantInformationReq) (res *information.MerchantInformationRes, err error) {
	res = &information.MerchantInformationRes{}
	// ZoneList
	var zoneList []string
	for _, zoneDir := range zoneDirs {
		subZoneList := ReadFile(zoneDir, "")
		for _, subZone := range subZoneList {
			if strings.Compare("+VERSION", subZone) != 0 {
				zoneList = append(zoneList, subZone)
			}
		}
	}
	res.SupportTimeZone = zoneList

	var supportCurrencys []*information.SupportCurrency
	supportCurrencys = append(supportCurrencys, &information.SupportCurrency{
		Currency: "EUR",
		Symbol:   "€",
		Scale:    100,
	})
	supportCurrencys = append(supportCurrencys, &information.SupportCurrency{
		Currency: "USD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrencys = append(supportCurrencys, &information.SupportCurrency{
		Currency: "JPY",
		Symbol:   "¥",
		Scale:    1,
	})
	res.SupportCurrency = supportCurrencys

	return res, nil
}

var zoneDirs = []string{
	// Update path according to your OS
	"/usr/share/zoneinfo/",
	"/usr/share/lib/zoneinfo/",
	"/usr/lib/locale/TZ/",
}

func ReadFile(zoneDir string, path string) []string {
	var zoneList []string
	files, _ := os.ReadDir(zoneDir + path)
	for _, f := range files {
		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
			continue
		}
		if f.IsDir() {
			subZoneList := ReadFile(zoneDir, path+"/"+f.Name())
			for _, subZone := range subZoneList {
				zoneList = append(zoneList, subZone)
			}
		} else {
			zoneList = append(zoneList, (path + "/" + f.Name())[1:])
		}
	}
	return zoneList
}
