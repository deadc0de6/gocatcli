package colorme

import (
	"fmt"

	"github.com/TwiN/go-color"
)

var (
	// UseColors enables the user of colors
	UseColors = true
)

// ColorMe object
type ColorMe struct {
	inline bool
}

// InUnderline underline
func (c *ColorMe) InUnderline(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return txt
	}
	return color.InUnderline(txt)
}

// InGray in gray
func (c *ColorMe) InGray(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return fmt.Sprintf("[gray]%s[-]", txt)
	}
	return color.InGray(txt)
}

// InRed in red
func (c *ColorMe) InRed(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return fmt.Sprintf("[red]%s[-]", txt)
	}
	return color.InRed(txt)
}

// InBlue in blue
func (c *ColorMe) InBlue(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return fmt.Sprintf("[blue]%s[-]", txt)
	}
	return color.InBlue(txt)
}

// InYellow in yellow
func (c *ColorMe) InYellow(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return fmt.Sprintf("[yellow]%s[-]", txt)
	}
	return color.InYellow(txt)
}

// InPurple in purple
func (c *ColorMe) InPurple(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return fmt.Sprintf("[purple]%s[-]", txt)
	}
	return color.InPurple(txt)
}

// InGreen in green
func (c *ColorMe) InGreen(txt string) string {
	if !UseColors {
		return txt
	}
	if c.inline {
		return fmt.Sprintf("[green]%s[-]", txt)
	}
	return color.InGreen(txt)
}

// NewColorme creates a new object
func NewColorme(inline bool) *ColorMe {
	cm := &ColorMe{
		inline: inline,
	}
	return cm
}
