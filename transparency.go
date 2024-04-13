package gopdf

import (
	"fmt"
	"sync"

	"errors"
)

type BlendModeType string

const (
	Hue             BlendModeType = "/Hue"
	Color           BlendModeType = "/Color"
	NormalBlendMode BlendModeType = "/Normal"
	Darken          BlendModeType = "/Darken"
	Screen          BlendModeType = "/Screen"
	Overlay         BlendModeType = "/Overlay"
	Lighten         BlendModeType = "/Lighten"
	Multiply        BlendModeType = "/Multiply"
	Exclusion       BlendModeType = "/Exclusion"
	ColorBurn       BlendModeType = "/ColorBurn"
	HardLight       BlendModeType = "/HardLight"
	SoftLight       BlendModeType = "/SoftLight"
	Difference      BlendModeType = "/Difference"
	Saturation      BlendModeType = "/Saturation"
	Luminosity      BlendModeType = "/Luminosity"
	ColorDodge      BlendModeType = "/ColorDodge"
)

const DefaultAplhaValue = 1

// Transparency defines an object alpha.
type Transparency struct {
	extGStateIndex int
	Alpha          float64
	BlendModeType  BlendModeType
}

func NewTransparency(alpha float64, blendModeType string) (Transparency, error) {
	if alpha < 0.0 || alpha > 1.0 {
		return Transparency{}, fmt.Errorf("alpha value is out of range (0.0 - 1.0): %.3f", alpha)
	}

	bmtType, err := defineBlendModeType(blendModeType)
	if err != nil {
		return Transparency{}, err
	}

	return Transparency{
		Alpha:         alpha,
		BlendModeType: bmtType,
	}, nil
}

func (t Transparency) GetId() string {
	keyStr := fmt.Sprintf("%.3f_%s", t.Alpha, t.BlendModeType)

	return keyStr
}

type TransparencyMap struct {
	syncer sync.Mutex
	table  map[string]Transparency
}

func NewTransparencyMap() TransparencyMap {
	return TransparencyMap{
		syncer: sync.Mutex{},
		table:  make(map[string]Transparency),
	}
}

func (tm *TransparencyMap) Find(transparency Transparency) (Transparency, bool) {
	key := transparency.GetId()

	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	t, ok := tm.table[key]
	if !ok {
		return Transparency{}, false
	}

	return t, ok

}

func (tm *TransparencyMap) Save(transparency Transparency) Transparency {
	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	key := transparency.GetId()
	tm.table[key] = transparency

	return transparency
}

func defineBlendModeType(bmType string) (BlendModeType, error) {
	switch bmType {
	case string(Hue):
		return Hue, nil
	case string(Color):
		return Color, nil
	case "", string(NormalBlendMode):
		return NormalBlendMode, nil
	case string(Darken):
		return Darken, nil
	case string(Screen):
		return Screen, nil
	case string(Overlay):
		return Overlay, nil
	case string(Lighten):
		return Lighten, nil
	case string(Multiply):
		return Multiply, nil
	case string(Exclusion):
		return Exclusion, nil
	case string(ColorBurn):
		return ColorBurn, nil
	case string(HardLight):
		return HardLight, nil
	case string(SoftLight):
		return SoftLight, nil
	case string(Difference):
		return Difference, nil
	case string(Saturation):
		return Saturation, nil
	case string(Luminosity):
		return Luminosity, nil
	case string(ColorDodge):
		return ColorDodge, nil
	default:
		return "", errors.New("blend mode is unknown")
	}
}
