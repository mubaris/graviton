package graviton

import (
	"log"
	"runtime"

	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sourcegraph/go-webkit2/webkit2"
)

// App contains metadata about the Graviton application.
type App struct {
	name          string
	width, height int
}

// Driver implements the Graviton application.
type Driver struct {
	App
	window *gtk.Window
	view   *webkit2.WebView
}

// NewApp creates and returns the application.
func NewApp(name string, width int, height int) *App {
	return &App{
		name:   name,
		width:  width,
		height: height,
	}
}

// NewDriver creates the Gtk application and WebView.
func NewDriver(app App) *Driver {
	runtime.LockOSThread()
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window: ", err)
	}
	win.SetDefaultSize(app.width, app.width)
	win.SetTitle(app.name)
	v := webkit2.NewWebView()
	return &Driver{
		App:    app,
		window: win,
		view:   v,
	}
}

// Initialize method initializes the application.
func (driver *Driver) Initialize() {
	driver.view.Connect("load-changed", func(_ *glib.Object, i int) {
		loadEvent := webkit2.LoadEvent(i)
		switch loadEvent {
		case webkit2.LoadFinished:
			fmt.Println("Load Finished")
		}
	})
	driver.window.Connect("destroy", func() {
		gtk.MainQuit()
		driver.view.Destroy()
	})
}

// Start method starts the application.
func (driver *Driver) Start() {
	driver.window.Add(driver.view)
	driver.window.ShowAll()
	gtk.Main()
}

// AttachURI method loads the URI.
func (driver *Driver) AttachURI(uri string) {
	driver.view.LoadURI(uri)
}

// AttachHTML methods loads the HTML.
func (driver *Driver) AttachHTML(content string, baseURI string) {
	driver.view.LoadHTML(content, baseURI)
}
