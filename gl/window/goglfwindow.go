package window

import (
	"log"

	"github.com/dabasan/go-dh3dbasis/coloru8"
	"github.com/dabasan/go-dh3dbasis/vector"
	"github.com/dabasan/goglf/gl/wrapper"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/dabasan/goglf/gl/front"
)

type ReshapeFunc func(gw *GOGLFWindow, width int, height int)
type UpdateFunc func(gw *GOGLFWindow)
type DrawFunc func(gw *GOGLFWindow)

type GOGLFWindow struct {
	Window *glfw.Window

	key_caf          *keyCountsAndFlags
	mouse_button_caf *mouseButtonCountsAndFlags

	reshape_func ReshapeFunc
	update_func  UpdateFunc
	draw_func    DrawFunc

	Background_color coloru8.ColorU8

	User_data interface{}
}

func NewGOGLFWindow(width int, height int, title string) (*GOGLFWindow, error) {
	gw := new(GOGLFWindow)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}
	log.Printf("info: OpenGL version=%v", gl.GoStr(gl.GetString(gl.VERSION)))

	front.Initialize()

	gw.key_caf = newKeyCountsAndFlags()
	gw.mouse_button_caf = newMouseButtonCountsAndFlags()

	window.SetKeyCallback(gw.keyCallback)
	window.SetMouseButtonCallback(gw.mouseButtonCallback)
	window.SetFramebufferSizeCallback(gw.framebufferSizeCallback)
	gw.Window = window

	gw.reshape_func = Reshape
	gw.update_func = Update
	gw.draw_func = Draw

	gw.Background_color = coloru8.GetColorU8FromFloat32Components(0.0, 0.0, 0.0, 1.0)

	return gw, nil
}

func (gw *GOGLFWindow) keyCallback(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	switch a {
	case glfw.Press:
		gw.key_caf.pressing_flags[k] = true
	case glfw.Release:
		gw.key_caf.pressing_flags[k] = false
	}
}
func (gw *GOGLFWindow) mouseButtonCallback(w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey) {
	switch a {
	case glfw.Press:
		gw.mouse_button_caf.pressing_flags[b] = true
	case glfw.Release:
		gw.mouse_button_caf.pressing_flags[b] = false
	}
}
func (gw *GOGLFWindow) framebufferSizeCallback(w *glfw.Window, width int, height int) {
	front.UpdateCameraAspect(width, height)
	gw.reshape_func(gw, width, height)
}

func (gw *GOGLFWindow) ClearDrawScreen() {
	wrapper.ClearColor(gw.Background_color.R, gw.Background_color.G, gw.Background_color.B, gw.Background_color.A)
	wrapper.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (gw *GOGLFWindow) display() {
	gw.updateKeyCounts()
	gw.updateMouseButtonCounts()

	gw.ClearDrawScreen()
	front.UpdateLighting()
	front.UpdateFog()

	gw.update_func(gw)
	front.UpdateCamera()
	gw.draw_func(gw)
}
func (gw *GOGLFWindow) updateKeyCounts() {
	for key, val := range gw.key_caf.pressing_flags {
		if _, ok := gw.key_caf.pressing_counts[key]; !ok {
			gw.key_caf.pressing_counts[key] = 0
		}
		if _, ok := gw.key_caf.releasing_counts[key]; !ok {
			gw.key_caf.releasing_counts[key] = 0
		}

		if val == true {
			gw.key_caf.pressing_counts[key]++
			gw.key_caf.releasing_counts[key] = 0
		} else {
			gw.key_caf.releasing_counts[key]++
			gw.key_caf.pressing_counts[key] = 0
		}
	}
}
func (gw *GOGLFWindow) updateMouseButtonCounts() {
	for key, val := range gw.mouse_button_caf.pressing_flags {
		if _, ok := gw.mouse_button_caf.pressing_counts[key]; !ok {
			gw.mouse_button_caf.pressing_counts[key] = 0
		}
		if _, ok := gw.mouse_button_caf.releasing_counts[key]; !ok {
			gw.mouse_button_caf.releasing_counts[key] = 0
		}

		if val == true {
			gw.mouse_button_caf.pressing_counts[key]++
			gw.mouse_button_caf.releasing_counts[key] = 0
		} else {
			gw.mouse_button_caf.releasing_counts[key]++
			gw.mouse_button_caf.pressing_counts[key] = 0
		}
	}
}

func Reshape(gw *GOGLFWindow, width int, height int) {

}
func Update(gw *GOGLFWindow) {
	front.SetCameraPositionAndTarget_UpVecY(
		vector.VGet(50.0, 50.0, 50.0), vector.VGet(0.0, 0.0, 0.0))
}
func Draw(gw *GOGLFWindow) {

}
func OnWindowClosing(gw *GOGLFWindow) {

}

func (gw *GOGLFWindow) InLoop() {
	gw.Window.MakeContextCurrent()
	gw.display()
	gw.Window.SwapBuffers()
}

func (gw *GOGLFWindow) SetReshapeFunc(f ReshapeFunc) {
	gw.reshape_func = f
}
func (gw *GOGLFWindow) SetUpdateFunc(f UpdateFunc) {
	gw.update_func = f
}
func (gw *GOGLFWindow) SetDrawFunc(f DrawFunc) {
	gw.draw_func = f
}

func (gw *GOGLFWindow) GetKeyPressingCount(k glfw.Key) int {
	val, ok := gw.key_caf.pressing_counts[k]
	if ok == false {
		return 0
	} else {
		return val
	}
}
func (gw *GOGLFWindow) GetkeyReleasingCount(k glfw.Key) int {
	val, ok := gw.key_caf.releasing_counts[k]
	if ok == false {
		return 0
	} else {
		return val
	}
}
func (gw *GOGLFWindow) GetMousePressingCount(b glfw.MouseButton) int {
	val, ok := gw.mouse_button_caf.pressing_counts[b]
	if ok == false {
		return 0
	} else {
		return val
	}
}
func (gw *GOGLFWindow) GetMouseReleasingCount(b glfw.MouseButton) int {
	val, ok := gw.mouse_button_caf.releasing_counts[b]
	if ok == false {
		return 0
	} else {
		return val
	}
}
