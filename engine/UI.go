package engine

import (
	"time"

	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
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
	//右侧装备栏
	rightPanel, _ := gui.NewImage("asset/UI/sidepanel_r.png")
	rightPanel.SetPosition(900-730/4.5, 0)
	rightPanel.SetSize(730/4.5, 2160/4.5)
	u.engine.Scence.Add(rightPanel)
	eq, _ := gui.NewImage("asset/UI/background.png")
	eq.SetPosition(850-1162/4.5, 47)
	eq.SetSize(1162/4.5, 1507/4.5)
	u.engine.Scence.Add(eq)
	//HP 动画
	tex1, _ := texture.NewTexture2DFromImage("asset/UI/HP.png")
	hp := gui.NewImageFromTex(tex1)
	anim1 := texture.NewAnimator(tex1, 46, 1)
	anim1.SetDispTime(16666 * time.Microsecond)
	u.anims = append(u.anims, anim1)
	//MP 动画
	tex1, _ = texture.NewTexture2DFromImage("asset/UI/MP.png")
	MP := gui.NewImageFromTex(tex1)
	anim2 := texture.NewAnimator(tex1, 46, 1)
	anim2.SetDispTime(16666 * time.Microsecond)
	u.anims = append(u.anims, anim2)
	//HP 缩放比例
	hpwd := 900 / 109.0
	hpdd := 900 / 86.0
	//MP 缩放比例
	mpwd := 900 / 109.0
	mpdd := 900 / 715.0
	u.engine.Scence.Add(hp)
	u.engine.Scence.Add(MP)

	//UI
	UIPanle, _ := gui.NewImage("asset/UI/front_panel.png")
	u.engine.Scence.Add(UIPanle)

	//窗体大小变化监听
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := u.engine.app.GetSize()
		u.engine.app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		u.engine.camera.SetAspect(float32(width) / float32(height))
		//更新UI大小尺寸
		UIPanle.SetSize(float32(width), float32(width)/6.8)
		UIPanle.SetPosition(0, float32(height)-float32(width)/6.8)
		//更新HP
		hp.SetSize(float32(width)/float32(hpwd), float32(width)/float32(hpwd))
		hp.SetPosition(float32(width)/float32(hpdd), float32(height)-float32(width)/float32(hpwd))
		//更新MP
		MP.SetSize(float32(width)/float32(mpwd), float32(width)/float32(mpwd))
		MP.SetPosition(float32(width)/float32(mpdd), float32(height)-float32(width)/float32(mpwd))
	}
	u.engine.app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	logo, _ := gui.NewImage("asset/UI/logo.png")
	logo.SetPosition(0, 0)
	logo.SetSize(648/2.0, 570/2.0)
	u.engine.Scence.Add(logo)

	// Creates the raycaster
	u.engine.rc = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	u.engine.rc.LinePrecision = 0.05
	u.engine.rc.PointPrecision = 0.05
	u.engine.app.Subscribe(window.OnMouseDown, u.onMouse)

	//fonts
	//fonts, _ := text.NewFont("asset/font/DiabloLight.ttf")
	// b1 := gui.NewLabel("diablo demo")
	// b1.SetFontSize(50)
	// b1.SetColor(&math32.Color{R: 1, G: 0, B: 0})
	// b1.SetFont(fonts)
	// b1.SetPosition(400, 0)
	// u.engine.Scence.Add(b1)

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
