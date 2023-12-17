package main

import (
	"scene_graph/obj"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

func main() {

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create a blue Cube
	mesh1 := graphic.NewMesh(
		geometry.NewCube(1),
		material.NewStandard(math32.NewColor("DarkBlue")),
	)
	mesh1.SetPosition(0, 0, 0)
	scene.Add(mesh1)

	// Create a red Cube
	mesh2 := graphic.NewMesh(
		geometry.NewCube(0.5),
		material.NewStandard(math32.NewColor("DarkRed")),
	)
	mesh2.SetPosition(1.5, 0, 0)
	scene.Add(mesh2)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	objScene := obj.NewScene()

	obj1 := objScene.NewObject()
	mesh1.WorldPosition(&obj1.Pos)
	objScene.Add(obj1)

	obj2 := objScene.NewObject()
	mesh2.WorldPosition(&obj2.Pos)
	obj1.Add(obj2)

	objScene.UpdateMatrixWorld(false)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		obj1.RotateY(float32(deltaTime.Seconds()))
		obj2.RotateZ(float32(deltaTime.Seconds()))

		objScene.UpdateMatrixWorld(false)

		// Update mesh position
		mesh1.SetPositionVec(obj1.GetWorldPosition())
		mesh2.SetPositionVec(obj2.GetWorldPosition())
		// Update mesh rotation
		mesh1.SetRotationQuat(obj1.GetWorldQuaternion())
		mesh2.SetRotationQuat(obj2.GetWorldQuaternion())

		renderer.Render(scene, cam)
	})
}
