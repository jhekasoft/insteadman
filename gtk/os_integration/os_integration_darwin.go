// +build darwin

package os_integration

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
Integrate(void) {
    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

    id menubar = [[NSMenu new] autorelease];
    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    id appMenu = [[NSMenu new] autorelease];
    id quitTitle = @"Quit InsteadMan";
    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
        action:@selector(terminate:) keyEquivalent:@"q"]
          	autorelease];
    [appMenu addItem:quitMenuItem];
    [appMenuItem setSubmenu:appMenu];
    [NSApp setMainMenu:menubar];

	[NSApp finishLaunching];
    [NSApp activateIgnoringOtherApps:YES];
    return 0;
}
*/
import "C"

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func OsIntegrate() {
	C.Integrate()
}

func OsIntegrateWindow(win *gtk.Window) {
	win.Connect("key-press-event", func(s *gtk.Window, ev *gdk.Event) {
		eventKey := &gdk.EventKey{ev}

		if isCmdQPresses(eventKey) {
			gtk.MainQuit()
		}
	})
}

func OsIntegrateDialog(win *gtk.Dialog) {
	win.Connect("key-press-event", func(s *gtk.Dialog, ev *gdk.Event) {
		eventKey := &gdk.EventKey{ev}

		if isCmdQPresses(eventKey) {
			s.Response(gtk.RESPONSE_CLOSE)
		}
	})
}

func isCmdQPresses(eventKey *gdk.EventKey) bool {
	// Check is cmd+q (or cmd+й for cyrillic) pressed
	return eventKey.State()&gdk.GDK_MOD2_MASK != 0 &&
		(eventKey.KeyVal() == gdk.KEY_q || eventKey.KeyVal() == gdk.KEY_Cyrillic_shorti)
}
