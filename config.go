package gopdf

// The units that can be used in the document
const (
	Unit_Unset = iota // No units were set, when conversion is called on nothing will happen
	Unit_PT           // Points
	Unit_MM           // Millimeters
	Unit_CM           // centimeters
	Unit_IN           // inches

	// The math needed to convert units to points
	conversion_Unit_PT = 1.0
	conversion_Unit_MM = 72.0 / 25.4
	conversion_Unit_CM = 72.0 / 2.54
	conversion_Unit_IN = 72.0
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
	case Unit_PT:
		return u * conversion_Unit_PT
	case Unit_MM:
		return u * conversion_Unit_MM
	case Unit_CM:
		return u * conversion_Unit_CM
	case Unit_IN:
		return u * conversion_Unit_IN
	default:
		return u
	}
}

// PointsToUnits converts points to the provided units
func PointsToUnits(t int, u float64) float64 {
	switch t {
	case Unit_PT:
		return u / conversion_Unit_PT
	case Unit_MM:
		return u / conversion_Unit_MM
	case Unit_CM:
		return u / conversion_Unit_CM
	case Unit_IN:
		return u / conversion_Unit_IN
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
