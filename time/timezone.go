package time

import (
	"os"
	"strings"
)

func CheckTimeZone(targetTimeZone string) bool {
	if len(targetTimeZone) > 0 {
		return false
	}
	timeZones := GetTimeZoneList()
	for _, zone := range timeZones {
		if strings.Compare(zone, targetTimeZone) == 0 {
			return true
		}
	}
	return false
}

func GetTimeZoneList() []string {
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
	return zoneList
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
