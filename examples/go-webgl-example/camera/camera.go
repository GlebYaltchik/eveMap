package camera

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/soypat/three"
)

type OrbitControl struct {
	o *js.Object
}

func NewOrbitControl(c three.PerspectiveCamera) *OrbitControl {
	return &OrbitControl{
		o: ocFunc.New(c.Object),
	}
}

var (
	ocFunc = js.Global.Get("OrbitControl").Invoke(js.Global.Get("THREE")) //nolint:gochecknoglobals
)
