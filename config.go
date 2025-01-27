package gopdf

// The units that can be used in the document
const (
	UnitUnset = iota // No units were set, when conversion is called on nothing will happen
	UnitPT           // Points
	UnitMM           // Millimeters
	UnitCM           // Centimeters
	UnitIN           // Inches
	UnitPX           // Pixels

	// The math needed to convert units to points
	conversionUnitPT = 1.0
	conversionUnitMM = 72.0 / 25.4
	conversionUnitCM = 72.0 / 2.54
	conversionUnitIN = 72.0
	//We use a dpi of 96 dpi as the default, so we get a conversionUnitPX = 3.0 / 4.0, which comes from 72.0 / 96.0.
	//If you want to change this value, you can change it at Config.ConversionForUnit
	//example: If you use dpi at 300.0
	//pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4, ConversionForUnit: 72.0 / 300.0 })
	conversionUnitPX = 3.0 / 4.0
)

// The units that can be used in the document (for backward compatibility)
// Deprecated: Use UnitUnset,UnitPT,UnitMM,UnitCM,UnitIN  instead
const (
	Unit_Unset = UnitUnset // No units were set, when conversion is called on nothing will happen
	Unit_PT    = UnitPT    // Points
	Unit_MM    = UnitMM    // Millimeters
	Unit_CM    = UnitCM    // Centimeters
	Unit_IN    = UnitIN    // Inches
	Unit_PX    = UnitPX    // Pixels
)

// Config static config
type Config struct {
	Unit int // The unit type to use when composing the document.
	//Value that use to convert units to points.
	//If this variable is not 0. This value will be used to calculate the unit conversion instead of the existing const value in the system.
	//And if this variable is not 0. Value ​​in Config.Unit will not be used.
	ConversionForUnit float64
	TrimBox           Box                 // The default trim box for all pages in the document
	PageSize          Rect                // The default page size for all pages in the document
	K                 float64             // Not sure
	Protection        PDFProtectionConfig // Protection settings
}

func (c Config) getUnit() int {
	return c.Unit
}
func (c Config) getConversionForUnit() float64 {
	return c.ConversionForUnit
}

// PDFProtectionConfig config of pdf protection
type PDFProtectionConfig struct {
	UseProtection bool
	Permissions   int
	UserPass      []byte
	OwnerPass     []byte
}

// UnitsToPoints converts units of the provided type to points
func UnitsToPoints(t int, u float64) float64 {
	return unitsToPoints(defaultUnitConfig{Unit: t}, u)
}

func unitsToPoints(unitCfg unitConfigurator, u float64) float64 {
	if unitCfg.getConversionForUnit() != 0 {
		return u * unitCfg.getConversionForUnit()
	}
	switch unitCfg.getUnit() {
	case UnitPT:
		return u * conversionUnitPT
	case UnitMM:
		return u * conversionUnitMM
	case UnitCM:
		return u * conversionUnitCM
	case UnitIN:
		return u * conversionUnitIN
	case UnitPX:
		return u * conversionUnitPX
	default:
		return u
	}
}

// PointsToUnits converts points to the provided units
func PointsToUnits(t int, u float64) float64 {
	return pointsToUnits(defaultUnitConfig{Unit: t}, u)
}

func pointsToUnits(unitCfg unitConfigurator, u float64) float64 {
	if unitCfg.getConversionForUnit() != 0 {
		return u / unitCfg.getConversionForUnit()
	}
	switch unitCfg.getUnit() {
	case UnitPT:
		return u / conversionUnitPT
	case UnitMM:
		return u / conversionUnitMM
	case UnitCM:
		return u / conversionUnitCM
	case UnitIN:
		return u / conversionUnitIN
	case UnitPX:
		return u / conversionUnitPX
	default:
		return u
	}
}

// UnitsToPointsVar converts units of the provided type to points for all variables supplied
func UnitsToPointsVar(t int, u ...*float64) {
	unitsToPointsVar(defaultUnitConfig{Unit: t}, u...)
}

func unitsToPointsVar(unitCfg unitConfigurator, u ...*float64) {
	for x := 0; x < len(u); x++ {
		*u[x] = unitsToPoints(unitCfg, *u[x])
	}
}

// PointsToUnitsVar converts points to the provided units for all variables supplied
func PointsToUnitsVar(t int, u ...*float64) {
	pointsToUnitsVar(defaultUnitConfig{Unit: t}, u...)
}

func pointsToUnitsVar(unitCfg unitConfigurator, u ...*float64) {
	for x := 0; x < len(u); x++ {
		*u[x] = pointsToUnits(unitCfg, *u[x])
	}
}

type unitConfigurator interface {
	getUnit() int
	getConversionForUnit() float64
}

type defaultUnitConfig struct {
	Unit              int
	ConversionForUnit float64
}

func (d defaultUnitConfig) getUnit() int {
	return d.Unit
}
func (d defaultUnitConfig) getConversionForUnit() float64 {
	return d.ConversionForUnit
}
