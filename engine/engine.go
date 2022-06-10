package engine

import (
	"fmt"
	"runtime"
	"time"

	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/audio"
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
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
)

var offset math32.Vector3

//Play state
const (
	ATTACK = iota
	IDLE
	WALK
)

//游戏类
type Game struct {
	camera     *camera.Camera
	app        *app.Application
	frameRater *util.FrameRater
	anims      []*animation.Animation
	anmisindex uint8
	Scence     *core.Node
	man        *core.Node
	rc         *collision.Raycaster
	target     math32.Vector3
	State      int
	audio      *Muscic
	ui         *UI
}

//游戏类实例化
func New() *Game {

	appli := app.App(900, 550, "golang game")
	//glfw.GetCurrentContext().SetSizeLimits(900, 550, 900, 550)
	cursorIcon, _ := appli.CreateCursor("asset/UI/mouse.png", 0, 0)
	appli.SetCursor(cursorIcon)
	//创建相机
	cam := camera.New(1)
	cam.SetQuaternion(-0.2640148, -0.4865232, -0.15705174, 0.817879)
	cam.SetPosition(-2.92947, 2.9979727, 1.4749823)

	//camera.NewOrbitControl(cam)
	offset.Add(math32.NewVector3(0, 0, 0))
	offset.Sub(math32.NewVector3(2.92947, -2.9979727, -1.4749823))
	//创建场景
	scene := core.NewNode()
	scene.Add(cam)
	//UI
	gui.Manager().Set(scene)
	g := &Game{
		State:      IDLE,
		Scence:     scene,
		camera:     cam,
		app:        appli,
		frameRater: util.NewFrameRater(30),
		anims:      make([]*animation.Animation, 0, 10),
	}
	//UI 加载
	g.ui = NewUI(g)
	g.ui.GUI()
	//场景加载
	go func() {
		g.LoadScence()
		runtime.GC()
	}()
	//背景色
	g.app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
	//音乐
	g.audio = NewMusic()
	// Create audio listener and add it to the current camera
	listener := audio.NewListener()
	cdir := g.camera.Direction()
	listener.SetDirectionVec(&cdir)
	g.camera.Add(listener)
	// Start the music!
	g.audio.musicGame.SetGain(10)
	g.audio.musicGame.Play()
	//显示debug信息
	//g.Debug()
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
	g.Scence.Add(helper.NewAxes(0.5))
	g.Scence.Add(helper.NewGrid(100, 1, math32.NewColor("green")))
	//
}

func (g *Game) update(renderer *renderer.Renderer, deltaTime time.Duration) {
	g.frameRater.Start()
	g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
	//相机渲染
	err := renderer.Render(g.Scence, g.camera)

	if err != nil {
		panic(err)
	}
	//
	if g.man != nil {
		manPos := g.man.Position()
		var target math32.Vector3
		target.Add(&manPos)
		target.Add(&offset)
		g.camera.SetPositionVec(&target)
		//控制角色
		g.ControllerMan(deltaTime)
		//状态机
		g.PlayAnimation(deltaTime)
	}
	//
	g.ui.Update()
	g.frameRater.Wait()
}

//控制角色
func (g *Game) ControllerMan(deltaTime time.Duration) {
	manPos := g.man.Position()
	dis := manPos.DistanceTo(&g.target)
	//move
	if dis > 0.1 {
		dist := 1.5 * float32(deltaTime.Seconds())
		// Get direction
		direction := g.target
		direction = *direction.Sub(&manPos)
		direction.Normalize()
		direction.MultiplyScalar(dist)
		// Get world position
		var position math32.Vector3
		g.man.WorldPosition(&position)
		position.Add(&direction)
		position.Y = 0

		var manDir math32.Vector3
		d := direction
		g.man.WorldDirection(&manDir)
		g.man.RotateY(manDir.AngleTo(&d) - 1.5)
		g.man.SetPositionVec(&position)
		//g.man.LookAt(math32.NewVector3(g.target.X, g.man.Position().Y, g.target.Z), math32.NewVector3(0, 1, 0))
		// g.man.RotateY(math.Pi / 2)
	} else if g.State != ATTACK {
		g.State = IDLE
		g.man.SetPositionVec(&g.target)
	}
}

//状态机
func (g *Game) PlayAnimation(deltaTime time.Duration) {
	//播放动画
	if len(g.anims) > 0 && int(g.anmisindex) < len(g.anims) {
		if g.State == ATTACK {
			g.anims[g.State].SetSpeed(2)
		}
		g.anims[g.State].Update(float32(deltaTime.Seconds()))

	}
}

//启动游戏
func (g *Game) Run() {
	g.app.Run(g.update)
}
