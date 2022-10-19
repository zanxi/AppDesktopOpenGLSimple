package main

// Imports necessary packages
import (
		"runtime"
        "syscall"
        "unsafe"
        "go-libs/w32"
		"go-libs/gl/v2.1/gl"
)

var GlInitialized bool

func init() {
    runtime.LockOSThread()
}

func MakeIntResource(id uint16) (*uint16) {
    return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func WndProc(hWnd w32.HWND, msg uint32, wParam, lParam uintptr) (uintptr) {
    switch msg {
    case w32.WM_SIZE:
		// Make sure we have initialized the gl before we use it
		if GlInitialized == true {
			rc := w32.GetClientRect(hWnd)
			gl.Viewport(0, 0, rc.Right, rc.Bottom)
		}
		break
    case w32.WM_DESTROY:
        w32.PostQuitMessage(0)
        break
    default:
        return w32.DefWindowProc(hWnd, msg, wParam, lParam)
    }
    return 0
}

func main() {
    
    GlInitialized = false

	// Store this module handle
	hInstance := w32.GetModuleHandle("")

	// Get a UTF-16 class name string
	lpszClassName := syscall.StringToUTF16Ptr("GoOpenGL!--Class")

	var wcex w32.WNDCLASSEX
	wcex.Size  		= uint32(unsafe.Sizeof(wcex))
	wcex.Style  	= w32.CS_HREDRAW | w32.CS_VREDRAW
	wcex.WndProc   	= syscall.NewCallback(WndProc)
	wcex.ClsExtra   = 0
	wcex.WndExtra   = 0
	wcex.Instance   = hInstance
	wcex.Icon       = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))
	wcex.Cursor     = w32.LoadCursor(0, MakeIntResource(w32.IDC_ARROW))
	wcex.Background	= w32.COLOR_WINDOW + 11
	wcex.MenuName  	= nil
	wcex.ClassName 	= lpszClassName
	wcex.IconSm 	= w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))
    
    // Register the class
	w32.RegisterClassEx(&wcex)

	// Now create a window using the registred class name
	hWnd := w32.CreateWindowEx(
	0, lpszClassName, syscall.StringToUTF16Ptr("Go OpenGL!"), 
	w32.WS_OVERLAPPEDWINDOW | w32.WS_VISIBLE, 
	w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 400, 0, 0, hInstance, nil)
    
    // Get the window's device context
    hDC := w32.GetDC(hWnd)
    
    // Setup pixel format ( We use 32-bit pixel )
    
    var pfd w32.PIXELFORMATDESCRIPTOR
    pfd.Size        = uint16(unsafe.Sizeof(pfd))
    pfd.Version     = 1;
    pfd.DwFlags		= w32.PFD_DRAW_TO_WINDOW | w32.PFD_SUPPORT_OPENGL | w32.PFD_DOUBLEBUFFER;
    pfd.IPixelType	= w32.PFD_TYPE_RGBA;
    pfd.ColorBits	= 32;
    
    pf := w32.ChoosePixelFormat(hDC, &pfd)
    
    w32.SetPixelFormat(hDC, pf, &pfd)
    
    // Create OpenGL context
    hRC := w32.WglCreateContext(hDC)
    w32.WglMakeCurrent(hDC, hRC)
    
    // Init & Setup the OpenGL
    if err := gl.Init(); err != nil {
		panic(err)
	}
    
    // Flag as initialized
    GlInitialized = true
    
    // Enable depth test
    gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
    
    // Enable some lighting for a nice visual
	gl.Enable(gl.LIGHTING)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{1, 1, 1, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)
    
    // Show the window
	w32.ShowWindow(hWnd, w32.SW_SHOWDEFAULT)
	w32.UpdateWindow(hWnd)

	// The main loop
	var msg w32.MSG
    quit := false
    
    for {
        // Check for exit request
        if quit == true {
            break
        }
        // Find for window messages
        if w32.PeekMessage(&msg, 0, 0, 0, w32.PM_REMOVE) {
			if msg.Message == w32.WM_QUIT {
				quit = true
			}
			w32.TranslateMessage(&msg)
			w32.DispatchMessage(&msg)
		}
        // Render scene
        drawScene()
        w32.SwapBuffers(hDC)
    }
    
    // Delete the OpenGL context
    w32.WglMakeCurrent(w32.HDC(0), w32.HGLRC(0))
    w32.ReleaseDC(hWnd, hDC)
    w32.WglDeleteContext(hRC)
    
    // Also delete the main window we created
    w32.DestroyWindow(hWnd)
}

func drawScene() {
    // Clear the window first with gray color
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.ClearColor(0.5, 0.5, 0.5, 0.0)
    gl.ClearDepth(1)
    
    // Setup projection matrix
    gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
    
    // Setup model-view matrix
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
    
    // Slightly rotate the view for a 3D view
    gl.Translatef(0, 0, -3.0)
    gl.Rotatef(40.0, 1, 0, 0)
    gl.Rotatef(40.0, 0, 1, 0)
    
/*
	// Draws a simple gradient triangle
    gl.Begin(gl.TRIANGLES)
    gl.Color3f(1.0, 0.0, 0.0)
    gl.Vertex2i(0,  1)
    gl.Color3f(0.0, 1.0, 0.0)
    gl.Vertex2i(-1, -1)
    gl.Color3f(0.0, 0.0, 1.0)
    gl.Vertex2i(1, -1)
    gl.End()
*/
    // Draw a cube
    drawCube()
}

func drawCube() {
    gl.Begin(gl.QUADS)

	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1)

	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1)

	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)

	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)

	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)

	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)

	gl.End()
}





