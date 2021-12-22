package color

import (
	"github.com/soypat/three"

	"evemap/examples/go-webgl-example/eve"
)

//nolint:gochecknoglobals
var (
	highSec = newColoredMaterial("green")
	nullSec = newColoredMaterial("red")
	lowSec  = newColoredMaterial("orange")

	innerJump         = newLineColor("darkgrey")
	constellationJump = newLineColor("yellow")
	regionJump        = newLineColor("orange")
)

func BySecurityStatus(s float64) three.Material {
	if s <= 0 {
		return nullSec
	}

	if s > 0.5 {
		return highSec
	}

	return lowSec
}

func ByJumpType(j eve.JumpInfo) three.Material {
	if j.FromRegionID != j.ToRegionID {
		return regionJump
	}

	if j.FromConstellationID != j.ToConstellationID {
		return constellationJump
	}

	return innerJump
}

func newColoredMaterial(color string) *three.MeshBasicMaterial {
	params := three.NewMaterialParameters()
	params.Color = three.NewColor(color)

	return three.NewMeshBasicMaterial(params)
}

func newLineColor(color string) *three.LineBasicMaterial {
	lp := three.NewMaterialParameters()
	lp.Color = three.NewColor(color)

	return three.NewLineBasicMaterial(lp)
}
