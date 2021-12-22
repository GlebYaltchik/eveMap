package main

import (
	"sync"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/soypat/three"
	"honnef.co/go/js/dom"

	camera2 "evemap/examples/go-webgl-example/camera"
	"evemap/examples/go-webgl-example/color"
	"evemap/examples/go-webgl-example/console"
	"evemap/examples/go-webgl-example/eve"
)

func getSize(obj *js.Object) (w float64, h float64) {
	return obj.Get("clientWidth").Float(), obj.Get("clientHeight").Float()
}

func resizeRendererToDisplaySize(r *three.WebGLRenderer, canvas *js.Object) float64 {
	rWidth, rHeight := getSize(r.Get("domElement"))
	cWidth, cHeight := getSize(canvas)

	if rWidth == cWidth && rHeight == cHeight {
		return 0
	}

	r.SetSize(cWidth, cHeight, true)

	return cWidth / cHeight
}

func globalRegister(k string, v interface{}) {
	js.Global.Set(k, v)
}

func main() {
	l := console.New(dom.GetWindow())

	log := l.Log

	log("eve mapper staring")

	loadComplete := sync.WaitGroup{}
	loadComplete.Add(1)

	go func() {
		defer loadComplete.Done()

		eve.LoadSolarSystems(log)
		globalRegister("eveStars", eve.SolarSystems)

		eve.LoadJumps(log)
		globalRegister("eveJumps", eve.JumpsData)
	}()

	viewPort := dom.GetWindow().Document().GetElementByID("cube")

	width, height := getSize(viewPort.Underlying())

	renderer := three.NewWebGLRenderer()
	renderer.SetSize(width, height, true)

	viewPort.AppendChild(dom.WrapNode(renderer.Get("domElement")))

	camera := three.NewPerspectiveCamera(50, width/height, 1, 100000)
	camera.Position.Set(3000, 5000, 3000)

	log(camera2.NewOrbitControl(camera))

	js.Global.Set("eveCamera", camera)

	scene := three.NewScene()

	light2 := three.NewAmbientLight(three.NewColor("white"), 0.1)
	scene.Add(light2)

	stars := three.NewGroup()

	allJumps := make(map[int]*three.Line)
	globalRegister("eveJumpObjects", allJumps)

	go func() {
		loadComplete.Wait()

		for _, system := range eve.SolarSystemsByName {
			star := three.NewMesh(three.NewSphereGeometry(three.SphereGeometryParameters{
				Radius:         5,
				WidthSegments:  4,
				HeightSegments: 2,
				ThetaStart:     0,
				ThetaLength:    0,
				PhiStart:       0,
				PhiLength:      0,
			}), color.BySecurityStatus(system.Security))

			star.Position.Set(system.X, system.Y, system.Z)

			stars.Add(star)
		}

		for n, info := range eve.JumpsData {
			f := eve.SolarSystems[info.FromSolarSystemID]
			t := eve.SolarSystems[info.ToSolarSystemID]

			attr := three.NewBufferAttribute(make([]float32, 6), 3)
			attr.SetXYZ(0, f.X, f.Y, f.Z)
			attr.SetXYZ(1, t.X, t.Y, t.Z)
			attr.NeedsUpdate = true

			g := three.NewBufferGeometry()
			g.Call("addAttribute", "position", attr)

			l := three.NewLine(g, color.ByJumpType(info))

			allJumps[n] = l

			stars.Add(l)
		}

		scene.Add(stars)

		log("all stars placed on the map")

		camera.LookAt(0, 0, 0)
	}()

	// start animation
	var animate func(t time.Duration)
	animate = func(d time.Duration) {
		// dt := d.Seconds()

		if aspect := resizeRendererToDisplaySize(&renderer, viewPort.Underlying()); aspect != 0 {
			camera.Aspect = aspect
			camera.UpdateProjectionMatrix()
		}

		// stars.Rotation.Y = dt
		// stars.Rotation.Z = dt / 1.5
		// stars.Rotation.Y = -dt / 2

		// hue := int(dt * 36.0)
		// mp := three.NewMaterialParameters()
		// mp.Color = three.NewColor(fmt.Sprintf("hsl(%d, 100%%, 50%%)", hue))
		//
		// mat := three.NewLineBasicMaterial(mp)
		//
		// for _, j := range allJumps {
		// 	j.Material = mat
		// }
		//
		renderer.Render(scene, camera)

		dom.GetWindow().RequestAnimationFrame(animate)
	}

	dom.GetWindow().RequestAnimationFrame(animate)
}
