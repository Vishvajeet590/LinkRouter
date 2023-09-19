package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"os"
	"os/exec"
	"strings"
)

type Browsers struct {
	name []string
}

func main() {
	var ScreenMsg string

	browsers, err := DetectBrowsers()
	if err != nil {
		return
	}
	err = CheckDesktop()
	if err != nil {
		ScreenMsg = err.Error()
	} else {
		ScreenMsg = "App is Installed, you can set it as Default Browser"
	}

	args := os.Args
	if len(args) > 1 {
		println(args[1])
		a := app.New()
		w := a.NewWindow(args[1])

		grid := container.New(layout.NewGridLayout(3))
		for _, browser := range browsers.name {
			button := widget.NewButton(browser, caller(browser, args[1], a))
			w.SetContent(button)
			grid.Add(button)
		}
		w.SetContent(grid)
		w.ShowAndRun()
	} else {
		a := app.New()
		w := a.NewWindow("Link Router")
		w.Resize(fyne.NewSize(300, 200))
		w.SetContent(
			container.NewBorder(
				nil, // TOP of the container
				nil,
				nil, // Right
				nil, // Left
				container.NewCenter(
					widget.NewLabel(ScreenMsg),
				),
			),
		)
		w.ShowAndRun()
	}

}

func caller(browser, url string, app fyne.App) func() {
	return func() {
		RunBrowser(browser, url)
		app.Quit()
	}
}

func RunBrowser(browser, url string) error {
	switch browser {
	case "chrome":
		cmd := exec.Command("google-chrome", url, "&")
		if err := cmd.Run(); err != nil {
			return err
		}
	case "firefox":
		cmd := exec.Command("firefox", url, "&")
		if err := cmd.Run(); err != nil {
			return err
		}
	case "brave":
		cmd := exec.Command("brave-browser", url, "&")
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func DetectBrowsers() (*Browsers, error) {
	ls := exec.Command("ls", "/usr/share/applications")
	grep := exec.Command("grep", "-i", "-e", "brave", "-e", "firefox", "-e", "chrome")
	pipe, _ := ls.StdoutPipe()
	defer pipe.Close()
	grep.Stdin = pipe
	err := ls.Start()
	if err != nil {
		return nil, err
	}
	res, _ := grep.Output()
	browsers := strings.Split(string(res), "\n")
	if len(browsers) == 0 {
		return nil, fmt.Errorf("no browser Found")
	}

	b := Browsers{
		name: make([]string, 0),
	}
	for _, browser := range browsers {
		switch browser {
		case "firefox.desktop":
			b.name = append(b.name, "firefox")
		case "google-chrome.desktop":
			b.name = append(b.name, "chrome")
		case "brave-browser.desktop":
			b.name = append(b.name, "brave")
		}
	}

	return &b, nil
}

func CheckDesktop() error {
	if _, err := os.Stat("/usr/share/applications/LinkRouter.desktop"); err == nil {
		return nil
	} else {
		err := os.WriteFile("/usr/share/applications/LinkRouter.desktop", []byte("[Desktop Entry]\nEncoding=UTF-8\nVersion=1.0\nType=Application\nTerminal=false\nExec=/usr/bin/LinkRouter %u\nName=LinkRouter\nCategories=Network;WebBrowser;\nMimeType=text/html;text/xml;application/xhtml+xml;application/vnd.mozilla.xul+xml;text/mml;x-scheme-handler/http;x-scheme-handler/https;"), 0755)
		if err != nil {
			return UnableToWriteErr
		}
	}
	return nil
}
