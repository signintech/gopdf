package gopdf

// The units that can be used in the document
const (
	UnitUnset = iota // No units were set, when conversion is called on nothing will happen
	UnitPT           // Points
	UnitMM           // Millimeters
	UnitCM           // centimeters
	UnitIN           // inches

	// The math needed to convert units to points
	conversionUnitPT = 1.0
	conversionUnitMM = 72.0 / 25.4
	conversionUnitCM = 72.0 / 2.54
	conversionUnitIN = 72.0
)

// The units that can be used in the document (for backward compatibility)
// Deprecated: Use UnitUnset,UnitPT,UnitMM,UnitCM,UnitIN  instead
const (
	Unit_Unset = UnitUnset // No units were set, when conversion is called on nothing will happen
	Unit_PT    = UnitPT    // Points
	Unit_MM    = UnitMM    // Millimeters
	Unit_CM    = UnitCM    // centimeters
	Unit_IN    = UnitIN    // inches
)

//Config static config
type Config struct {
	Unit       int                 // The unit type to use when composing the document.
	PageSize   Rect                // The default page size for all pages in the document
	K          float64             // Not sure
	Protection PDFProtectionConfig // Protection settings
}

//PDFProtectionConfig config of pdf protection
type PDFProtectionConfig struct {
	UseProtection bool
	Permissions   int
	UserPass      []byte
	OwnerPass     []byte
}

// UnitsToPoints converts units of the provided type to points
func UnitsToPoints(t int, u float64) float64 {
	switch t {
	case UnitPT:
		return u * conversionUnitPT
	case UnitMM:
		return u * conversionUnitMM
	case UnitCM:
		return u * conversionUnitCM
	case UnitIN:
		return u * conversionUnitIN
	default:
		return u
	}
}

// PointsToUnits converts points to the provided units
func PointsToUnits(t int, u float64) float64 {
	switch t {
	case UnitPT:
		return u / conversionUnitPT
	case UnitMM:
		return u / conversionUnitMM
	case UnitCM:
		return u / conversionUnitCM
	case UnitIN:
		return u / conversionUnitIN
	default:
		return u
	}
}

// UnitsToPointsVar converts units of the provided type to points for all variables supplied
func UnitsToPointsVar(t int, u ...*float64) {
	for x := 0; x < len(u); x++ {
		*u[x] = UnitsToPoints(t, *u[x])
	}
}

// PointsToUnitsVar converts points to the provided units for all variables supplied
func PointsToUnitsVar(t int, u ...*float64) {
	for x := 0; x < len(u); x++ {
		*u[x] = PointsToUnits(t, *u[x])
	}
}
