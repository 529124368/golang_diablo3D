package engine

import (
	"fmt"
	"test/tools"
	"time"

	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/loader/gltf"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/window"
)

const (
	CMD_W = iota
	CMD_S
	CMD_A
	CMD_D
	CMD_LAST
)

var offset math32.Vector3

//游戏类
type Game struct {
	cam        *camera.Camera
	app        *app.Application
	frameRater *util.FrameRater
	anims      []*animation.Animation
	anmisindex uint8
	Scence     *core.Node
	stopAnm    bool
	man        *core.Node
	commad     [CMD_LAST]bool
	rc         *collision.Raycaster
}

//游戏类实例化
func New() *Game {
	appli := app.App(1200, 950, "golang game")

	//创建相机
	cam := camera.New(1)
	cam.SetQuaternion(-0.2640148, -0.4865232, -0.15705174, 0.817879)
	cam.SetPosition(-2.92947, 2.9979727, 1.4749823)
	fmt.Println(cam.Position())
	//camera.NewOrbitControl(cam)
	offset.Add(math32.NewVector3(0, 0, 0))
	offset.Sub(math32.NewVector3(2.92947, -2.9979727, -1.4749823))
	//创建场景
	scene := core.NewNode()
	scene.Add(cam)
	//UI
	gui.Manager().Set(scene)

	g := &Game{
		Scence:     scene,
		cam:        cam,
		app:        appli,
		frameRater: util.NewFrameRater(30),
		anims:      make([]*animation.Animation, 0, 10),
	}
	//UI 加载
	g.GUI()
	//场景加载
	g.LoadScence()
	//背景色
	g.app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
	//显示debug信息
	g.Debug()

	return g
}

//动画模型加载
func (g *Game) newPlayerModel(path string) {
	go func() {
		model, err := gltf.ParseJSON(path)
		if err != nil {
			fmt.Println(err)
		}
		gg, err := model.LoadScene(*model.Scene)
		g.man = gg.GetNode()
		if err != nil {
			fmt.Println(err)
		}

		g.man.SetScale(1, 1, 1)
		g.man.SetPosition(0, 0, 0)
		//加载动画
		if len(model.Animations) > 0 {
			for i := range model.Animations {
				anim, _ := model.LoadAnimation(i)
				anim.SetLoop(true)
				g.anims = append(g.anims, anim)
			}
		}
		g.Scence.Add(gg)
	}()
}

//普通模型加载
func (g *Game) newModel(path string, y float32) {
	go func() {
		model, err := gltf.ParseJSON(path)
		if err != nil {
			fmt.Println(err)
		}
		gg, _ := model.LoadScene(*model.Scene)
		gg.GetNode().SetScale(1, 1, 1)
		gg.GetNode().SetPosition(0, y, 0)
		g.Scence.Add(gg)
	}()
}

//删除人物
func (g *Game) DelPlayerModel() {
	g.Scence.Remove(g.man)
	g.anims = g.anims[:0]
}

//加载模型
func (g *Game) LoadScence() {

	//加载模型和动画
	g.newPlayerModel("asset/gltf/player/ba.gltf")

	//加载
	g.newModel("asset/gltf/player/am.gltf", 0)

	//加载地图
	g.newModel("asset/gltf/map/map.gltf", -0.85)

	//加载灯光
	//g.Scence.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	g.Scence.Add(pointLight)
	pointLight = light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(-2.5, 0, 0)
	g.Scence.Add(pointLight)
	pointLight = light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(-4.3, 0, 3.3)
	g.Scence.Add(pointLight)
	pointLight = light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(-5, 0, 2)
	g.Scence.Add(pointLight)
	pointLight1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 1.0)
	pointLight1.SetPosition(0, 10, 10)
	g.Scence.Add(pointLight1)

	//skybox
	skybox, err := graphic.NewSkybox(graphic.SkyboxData{
		DirAndPrefix: "asset/skybox/", Extension: "jpg",
		Suffixes: [6]string{"dark-s_nx", "dark-s_ny", "dark-s_nz", "dark-s_px", "dark-s_py", "dark-s_pz"}})
	if err != nil {
		panic(err)
	}
	g.Scence.Add(skybox)
}

//debug 显示
func (g *Game) Debug() {
	// // Create and add an axis helper to the scene
	//辅助显示
	// g.Scence.Add(helper.NewAxes(0.5))
	// g.Scence.Add(helper.NewGrid(100, 1, math32.NewColor("green")))
	//
}

//UI
func (g *Game) GUI() {

	selectF := tools.NewFileSelectButton("./", "Select File", 400, 300)
	selectF.SetPosition(200, 10)
	selectF.Subscribe("OnSelect", func(evname string, ev interface{}) {
		fpath := ev.(string)
		fmt.Println(fpath)
		g.DelPlayerModel()
		g.newPlayerModel(fpath)
	})
	g.Scence.Add(selectF)

	//注册监听
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := g.app.GetSize()
		g.app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		g.cam.SetAspect(float32(width) / float32(height))
	}
	g.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	//GUI
	fonts, _ := text.NewFont("asset/font/DiabloLight.ttf")

	//停止动画按钮
	btn := gui.NewButton("stop")
	btn.Label.SetFont(fonts)
	btn.SetPosition(66, 141)
	btn.SetSize(40, 40)

	//监听
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.stopAnm = !g.stopAnm
		g.anims[g.anmisindex].SetPaused(g.stopAnm)
	})
	g.Scence.Add(btn)

	//切换动画按钮
	btn1 := gui.NewButton("change 0")
	btn1.Label.SetFont(fonts)
	btn1.SetPosition(66, 191)
	btn1.SetSize(40, 40)

	//监听
	btn1.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.anmisindex = 0
	})
	g.Scence.Add(btn1)

	//切换动画按钮
	btn2 := gui.NewButton("change 1")
	btn2.Label.SetFont(fonts)
	btn2.SetPosition(66, 241)
	btn2.SetSize(40, 40)

	//监听
	btn2.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.anmisindex = 1
	})
	g.Scence.Add(btn2)

	//切换动画按钮
	btn3 := gui.NewButton("change 2")
	btn3.Label.SetFont(fonts)
	btn3.SetPosition(66, 291)
	btn3.SetSize(40, 40)

	//监听
	btn3.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.anmisindex = 2
	})
	g.Scence.Add(btn3)

	b1 := gui.NewLabel("diablo demo")
	b1.SetFontSize(50)
	b1.SetColor(&math32.Color{R: 1, G: 0, B: 0})
	b1.SetFont(fonts)
	b1.SetPosition(400, 141)
	g.Scence.Add(b1)

	//按键监听
	g.app.Subscribe(window.OnKeyDown, g.onKey)
	g.app.Subscribe(window.OnKeyUp, g.onKey)
	//
	// Creates the raycaster
	g.rc = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	g.rc.LinePrecision = 0.05
	g.rc.PointPrecision = 0.05
	g.app.Subscribe(window.OnMouseDown, g.onMouse)
}
func (g *Game) onMouse(evname string, ev interface{}) {
	// Convert mouse coordinates to normalized device coordinates
	mev := ev.(*window.MouseEvent)
	width, height := g.app.GetSize()
	x := 2*(mev.Xpos/float32(width)) - 1
	y := -2*(mev.Ypos/float32(height)) + 1

	// Set the raycaster from the current camera and mouse coordinates
	g.rc.SetFromCamera(g.cam, x, y)
	//fmt.Printf("rc:%+v\n", t.rc.Ray)

	// Checks intersection with all objects in the scene
	intersects := g.rc.IntersectObjects(g.Scence.Children(), true)
	//fmt.Printf("intersects:%+v\n", intersects)
	if len(intersects) == 0 {
		return
	}

	// Get first intersection
	obj := intersects[0].Object
	// Convert INode to IGraphic
	_, ok := obj.(graphic.IGraphic)
	if !ok {
		fmt.Println("Not graphic")
		return
	}
	fmt.Println(obj.Position())
}

func (g *Game) update(renderer *renderer.Renderer, deltaTime time.Duration) {
	g.frameRater.Start()

	g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
	//相机渲染
	err := renderer.Render(g.Scence, g.cam)

	if err != nil {
		panic(err)
	}
	//控制角色
	g.ControllerMan(deltaTime)
	//
	if g.man != nil {
		manPos := g.man.Position()
		var target math32.Vector3
		target.Add(&manPos)
		target.Add(&offset)
		g.cam.SetPositionVec(&target)
	}
	//播放动画
	if len(g.anims) > 0 && int(g.anmisindex) < len(g.anims) {
		g.anims[g.anmisindex].Update(float32(deltaTime.Seconds()))
	}
	g.frameRater.Wait()
}

//控制角色
func (g *Game) ControllerMan(deltaTime time.Duration) {

	if g.commad[CMD_A] {
		g.man.RotateY(0.2)
	}
	if g.commad[CMD_D] {
		g.man.RotateY(-0.2)
	}
	if g.commad[CMD_W] || g.commad[CMD_S] {
		// Calculates the distance to move
		dist := 0.9 * float32(deltaTime.Seconds())
		// Get direction
		var quat math32.Quaternion
		g.man.WorldQuaternion(&quat)
		direction := math32.Vector3{X: 1, Y: 0, Z: 0}
		direction.ApplyQuaternion(&quat)
		direction.Normalize()
		direction.MultiplyScalar(dist)
		if g.commad[CMD_S] {
			direction.Negate()
		}
		// Get world position
		var position math32.Vector3
		g.man.WorldPosition(&position)
		position.Add(&direction)
		g.man.SetPositionVec(&position)
	}
}

// Process key events
func (g *Game) onKey(evname string, ev interface{}) {
	var state bool
	if evname == window.OnKeyDown {
		state = true
	} else {
		state = false
	}
	kev := ev.(*window.KeyEvent)
	switch kev.Key {
	case window.KeyW:
		g.commad[CMD_W] = state
	case window.KeyS:
		g.commad[CMD_S] = state
	case window.KeyA:
		g.commad[CMD_A] = state
	case window.KeyD:
		g.commad[CMD_D] = state
	}
}

//启动游戏
func (g *Game) Run() {
	g.app.Run(g.update)
}
