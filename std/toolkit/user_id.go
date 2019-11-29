package toolkit

import (
	std "github.com/PKUJohnson/solar/std"
)

func GetUserIDFromDBId(appType string, dbId int64) int64 {
	aflag := std.AppFlagFromName(appType)
	return dbId + std.AppFlagStart*int64(aflag)
}

func GetAppFlagFromUid(uid int64) std.AppFlag {
	i := uid / std.AppFlagStart
	if i < 0 || int64(i) > int64(std.MaxAppFlag) {
		return std.InvalidAppFlag
	}
	return std.AppFlag(i)
}

func GetAppTypeFromUid(uid int64) string {
	return GetAppFlagFromUid(uid).Name()
}

func GetDBIdFromUserID(uid int64) int64 {
	f := GetAppFlagFromUid(uid)
	return uid - int64(f)*std.AppFlagStart
}

func GetDBIdsFromUserIDs(uids []int64) []int64 {
	f := GetAppFlagFromUid(uids[0])
	Ids := []int64{}
	for _, userid := range uids {
		Ids = append(Ids, userid-int64(f)*std.AppFlagStart)
	}
	return Ids
}
