package gopdf

// BreakMode type for text break modes.
type BreakMode int

const (
	// BreakModeStrict causes the text-line to break immediately in case the current character would not fit into
	// the processed text-line. The separator (if provided) will be attached accordingly as a line suffix
	// to stay within the defined width.
	BreakModeStrict BreakMode = iota

	// BreakModeIndicatorSensitive will try to break the current line based on the last index of a provided
	// BreakIndicator. If no indicator sensitive break can be performed a strict break will be performed,
	// potentially working with the given separator as a suffix.
	BreakModeIndicatorSensitive
)

var (
	// DefaultBreakOption will cause the text to break mid-word without any separator suffixes.
	DefaultBreakOption = BreakOption{
		Mode:           BreakModeStrict,
		BreakIndicator: 0,
		Separator:      "",
	}
)

// BreakOption allows to configure the behavior of splitting or breaking larger texts via SplitTextWithOption.
type BreakOption struct {
	// Mode defines the mode which should be used
	Mode BreakMode
	// BreakIndicator is taken into account when using indicator sensitive mode to avoid mid-word line breaks
	BreakIndicator rune
	// Separator will act as a suffix for mid-word breaks when using strict mode
	Separator string
}

func (bo BreakOption) HasSeparator() bool {
	return bo.Separator != ""
}
