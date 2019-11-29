package std

const (
	invalidApp   = "-"
	SolarApp     = "solar"

	invalidAppID = 0
	solarAppId   = 9

	nonspecificResType    = "-"
	resUser               = "user"
	nonspecificResTypeID  = 0
)

const (
	//MaxRawID = 1*1000*1000*1000*1000 - 1 // 2^52 ＝ 4 503 599 627 370 496
	//MaxID      = 4499 * 1000 * 1000 * 1000 * 1000

	AppFlagStart = 1 * 1000 * 1000 * 1000 * 1000
	ResFlagStart = 100 * 1000 * 1000 * 1000 * 1000

	InvalidAppFlag = AppFlag(invalidAppID)
	SolarAppFlag   = AppFlag(solarAppId)
	MaxAppFlag     = SolarAppFlag

	NonspecificResFlag  = ResFlag(nonspecificResTypeID)
)

const (
	REMIND = 1
)

type AppFlag int64

func AppFlagFromName(name string) AppFlag {
	if name == SolarApp {
		return SolarAppFlag
	}
	return InvalidAppFlag
}

func (a AppFlag) Name() string {
	if a == SolarAppFlag {
		return SolarApp
	}

	return invalidApp
}

func (a AppFlag) Int64() int64 {
	return int64(a)
}

func (a AppFlag) Valid() bool {
	return a > 0 && a <= MaxAppFlag
}

type ResFlag int64

func ResFlagFromName(name string) ResFlag {
	return NonspecificResFlag
}

func (r ResFlag) Name() string {
	return nonspecificResType
}

func (r ResFlag) Int64() int64 {
	return int64(r)
}

func (r ResFlag) Valid() bool {
	return r > 0
}
