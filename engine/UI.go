package engine

import (
	"test/config"
	"time"

	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
)

var openBagSize float32

type UI struct {
	engine    *Game
	anims     []*texture.Animator
	mainPanel *gui.Panel
	eqPanel   *gui.Panel
}

func NewUI(g *Game) *UI {
	ui := new(UI)
	ui.engine = g
	ui.mainPanel = gui.NewPanel(0, 0)
	ui.eqPanel = gui.NewPanel(0, 0)
	// Show and enable demo panel
	ui.mainPanel.Add(ui.eqPanel)
	ui.eqPanel.SetVisible(false)
	g.Scence.Add(ui.mainPanel)
	return ui
}

func (u *UI) GUI() {

	//右侧装备栏
	datas, _ := asset.ReadFile("asset/UI/sidepanel_r.png")
	rightPanel, _ := gui.NewImagePlus(datas)
	u.eqPanel.Add(rightPanel)

	datas, _ = asset.ReadFile("asset/UI/background.png")
	eq, _ := gui.NewImagePlus(datas)
	u.eqPanel.Add(eq)

	datas, _ = asset.ReadFile("asset/UI/sidepanel_hinge_r.png")
	rightClose, _ := gui.NewImagePlus(datas)
	u.eqPanel.Add(rightClose)

	//HP 动画
	datas, _ = asset.ReadFile("asset/UI/HP.png")
	tex1, _ := texture.NewTexture2DFromImagePlius(datas)
	hp := gui.NewImageFromTex(tex1)
	anim1 := texture.NewAnimator(tex1, 46, 1)
	anim1.SetDispTime(16666 * time.Microsecond)
	u.anims = append(u.anims, anim1)
	u.mainPanel.Add(hp)

	//MP 动画
	datas, _ = asset.ReadFile("asset/UI/MP.png")
	tex1, _ = texture.NewTexture2DFromImagePlius(datas)
	MP := gui.NewImageFromTex(tex1)
	anim2 := texture.NewAnimator(tex1, 46, 1)
	anim2.SetDispTime(16666 * time.Microsecond)
	u.anims = append(u.anims, anim2)
	u.mainPanel.Add(MP)

	//UI
	datas, _ = asset.ReadFile("asset/UI/front_panel.png")
	UIPanle, _ := gui.NewImagePlus(datas)
	u.mainPanel.Add(UIPanle)

	//按钮
	datas, _ = asset.ReadFile("asset/UI/btn_normal.png")
	openBag, _ := gui.NewImagePlus(datas)
	openBag.Subscribe(window.OnMouseDown, func(s string, i interface{}) {
		datas, _ = asset.ReadFile("asset/UI/btn_pressed.png")
		openBag.SetImagePlus(datas)
		openBag.SetSize(openBagSize, openBagSize)
		u.eqPanel.SetVisible(!u.eqPanel.Visible())
	})
	openBag.Subscribe(window.OnMouseUp, func(s string, i interface{}) {
		datas, _ = asset.ReadFile("asset/UI/btn_normal.png")
		openBag.SetImagePlus(datas)
		openBag.SetSize(openBagSize, openBagSize)
	})
	u.mainPanel.Add(openBag)

	//窗体大小变化监听
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := u.engine.app.GetSize()
		u.engine.app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		u.engine.camera.SetAspect(float32(width) / float32(height))
		//
		u.mainPanel.SetSize(float32(width), float32(height))
		u.eqPanel.SetSize(float32(width), float32(height))
		//更新UI大小尺寸
		UIPanle.SetSize(float32(width), float32(width)/6.8)
		UIPanle.SetPosition(0, float32(height)-float32(width)/6.8)
		//更新HP
		hp.SetSize(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/109.0), float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/109.0))
		hp.SetPosition(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/86.0), float32(height)-float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/109.0))
		//更新MP
		MP.SetSize(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/109.0), float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/109.0))
		MP.SetPosition(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/715.0), float32(height)-float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/109.0))
		//更新装备栏
		rightPanel.SetPosition(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/717), 0)
		rightPanel.SetSize(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/182), float32(height)/(float32(config.DEFAULT_SCREEN_HEIGHT)/540))
		eq.SetPosition(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/549), float32(height)/(float32(config.DEFAULT_SCREEN_HEIGHT)/52))
		eq.SetSize(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/290), float32(height)/(float32(config.DEFAULT_SCREEN_HEIGHT)/376))
		rightClose.SetPosition(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/817), float32(height)/(float32(config.DEFAULT_SCREEN_HEIGHT)/220))
		rightClose.SetSize(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/83), float32(height)/(float32(config.DEFAULT_SCREEN_HEIGHT)/74))
		//
		openBag.SetPosition(float32(width)/(float32(config.DEFAULT_SCREEN_WIDTH)/435), float32(height)/(float32(config.DEFAULT_SCREEN_HEIGHT)/510))
		openBagSize = float32(height) / (float32(config.DEFAULT_SCREEN_HEIGHT) / 40)
		openBag.SetSize(openBagSize, openBagSize)

	}
	u.engine.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// logo, _ := gui.NewImage("engine/asset/UI/logo.png")
	// logo.SetPosition(0, 0)
	// logo.SetSize(648/2.0, 570/2.0)
	// u.mainPanel.Add(logo)

	// Creates the raycaster
	u.engine.rc = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	u.engine.rc.LinePrecision = 0.05
	u.engine.rc.PointPrecision = 0.05
	u.engine.app.Subscribe(window.OnMouseDown, u.onMouse)

	//u.engine.app.Subscribe(, u.onMouse)
	//fonts
	//fonts, _ := text.NewFont("engine/asset/font/DiabloLight.ttf")
	// b1 := gui.NewLabel("diablo demo")
	// b1.SetFontSize(50)
	// b1.SetColor(&math32.Color{R: 1, G: 0, B: 0})
	// b1.SetFont(fonts)
	// b1.SetPosition(400, 0)
	// u.mainPanel.Add(b1)

	// selectF := tools.NewFileSelectButton("./", "Select File", 400, 300)
	// selectF.SetPosition(200, 10)
	// selectF.Subscribe("OnSelect", func(evname string, ev interface{}) {
	// 	fpath := ev.(string)
	// 	fmt.Println(fpath)
	// 	g.DelPlayerModel()
	// 	g.newPlayerModel(fpath)
	// })
	// g.Scence.Add(selectF)
}

//鼠标事件
func (u *UI) onMouse(evname string, ev interface{}) {
	mev := ev.(*window.MouseEvent)
	if mev.Button == window.MouseButtonLeft {
		width, height := u.engine.app.GetSize()
		x := 2*(mev.Xpos/float32(width)) - 1
		y := -2*(mev.Ypos/float32(height)) + 1

		u.engine.rc.SetFromCamera(u.engine.camera, x, y)

		// Checks intersection with all objects in the scene
		intersects := u.engine.rc.IntersectObjects(u.engine.Scence.Children(), true)
		if len(intersects) == 0 {
			return
		}
		u.engine.State = WALK
		// Get first intersection
		u.engine.target = intersects[0].Point
	} else if mev.Button == window.MouseButtonRight {
		u.engine.State = ATTACK
		u.engine.anims[u.engine.State].Reset()
		u.engine.anims[u.engine.State].SetPaused(false)
		//声音控制
		u.engine.audio.skill.SetGain(20)
		u.engine.audio.skill.Play()
	} else if mev.Button == window.MouseButtonMiddle {
		u.engine.State = JUMP
		u.engine.anims[u.engine.State].Reset()
		u.engine.anims[u.engine.State].SetPaused(false)
		//声音控制
		u.engine.audio.skill.SetGain(20)
		u.engine.audio.skill.Play()
	}
}

// Update is called every frame.
func (u *UI) Update() {
	for _, anim := range u.anims {
		anim.Update(time.Now())
	}
}
