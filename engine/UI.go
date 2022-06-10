package engine

import (
	"time"

	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
)

type UI struct {
	engine *Game
	anims  []*texture.Animator
}

func NewUI(g *Game) *UI {
	ui := new(UI)
	ui.engine = g
	return ui
}

func (u *UI) GUI() {
	//fonts
	fonts, _ := text.NewFont("asset/font/DiabloLight.ttf")
	//HP 动画
	tex1, _ := texture.NewTexture2DFromImage("asset/UI/HP.png")
	hp := gui.NewImageFromTex(tex1)
	anim1 := texture.NewAnimator(tex1, 46, 1)
	anim1.SetDispTime(16666 * time.Microsecond)
	u.anims = append(u.anims, anim1)
	hpwd := 900 / 109.0
	hpdd := 900 / 86.0
	u.engine.Scence.Add(hp)
	//UI 血槽
	im, _ := gui.NewImage("asset/UI/front_panel.png")
	u.engine.Scence.Add(im)

	//注册监听
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := u.engine.app.GetSize()
		u.engine.app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		u.engine.camera.SetAspect(float32(width) / float32(height))
		//更新UI大小尺寸
		im.SetSize(float32(width), float32(width)/6.8)
		im.SetPosition(0, float32(height)-float32(width)/6.8)
		//更新HP
		hp.SetSize(float32(width)/float32(hpwd), float32(width)/float32(hpwd))
		hp.SetPosition(float32(width)/float32(hpdd), float32(height)-float32(width)/float32(hpwd))
	}
	u.engine.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// //监听
	// btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
	// 	g.stopAnm = !g.stopAnm
	// 	g.anims[g.anmisindex].SetPaused(g.stopAnm)
	// })
	// g.Scence.Add(btn)

	// //切换动画按钮
	// btn1 := gui.NewButton("change 0")
	// btn1.Label.SetFont(fonts)
	// btn1.SetPosition(66, 191)
	// btn1.SetSize(40, 40)

	// //监听
	// btn1.Subscribe(gui.OnClick, func(name string, ev interface{}) {
	// 	g.anmisindex = 0
	// })
	// g.Scence.Add(btn1)

	// //切换动画按钮
	// btn2 := gui.NewButton("change 1")
	// btn2.Label.SetFont(fonts)
	// btn2.SetPosition(66, 241)
	// btn2.SetSize(40, 40)

	// //监听
	// btn2.Subscribe(gui.OnClick, func(name string, ev interface{}) {
	// 	g.anmisindex = 1
	// })
	// g.Scence.Add(btn2)

	// //切换动画按钮
	// btn3 := gui.NewButton("change 2")
	// btn3.Label.SetFont(fonts)
	// btn3.SetPosition(66, 291)
	// btn3.SetSize(40, 40)

	// //监听
	// btn3.Subscribe(gui.OnClick, func(name string, ev interface{}) {
	// 	g.anmisindex = 2
	// })
	// g.Scence.Add(btn3)

	b1 := gui.NewLabel("diablo demo")
	b1.SetFontSize(50)
	b1.SetColor(&math32.Color{R: 1, G: 0, B: 0})
	b1.SetFont(fonts)
	b1.SetPosition(400, 0)
	u.engine.Scence.Add(b1)

	// Creates the raycaster
	u.engine.rc = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	u.engine.rc.LinePrecision = 0.05
	u.engine.rc.PointPrecision = 0.05
	u.engine.app.Subscribe(window.OnMouseDown, u.onMouse)

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
