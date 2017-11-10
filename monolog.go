package main

import (
	ui "github.com/deepakkamesh/termui"
	//"github.com/nsf/termbox-go"
	"net/http"
)

const (
	url           = `http://localhost`
	scriptingPORT = `:8888/ss`
	taskingPORT   = `:12121/ts`
	agentPORT     = `:8081/agent`
	health        = `/health`
)
const (
	on = iota
	off
)

var (
	logText      = `123`
	log          = ui.NewPar(logText)
	scripting    = ui.NewPar(`>> scripting-MS`)
	tasking      = ui.NewPar(`>> tasking-MS`)
	agent        = ui.NewPar(`>> agent-MS`)
	scriptingURL = url + scriptingPORT + health
	taskingURL   = url + taskingPORT + health
	agentURL     = url + agentPORT + health + `Check`

	ms = []*ui.Par{scripting, tasking}
)

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	service()
	draw := func(t int) {
		ui.Render(scripting, tasking, agent, log)
	}
	handleKbdUp(log)
	ui.Handle(`/sys/kbd/q`, func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle(`/timer/1s`, func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})
	ui.Loop()
}

func handleKbdUp(log *ui.Par) {
	ui.Handle(`/sys/kbd/<up>`, func(e ui.Event) {
		log.TextFgColor = ui.ColorGreen
	})
	ui.Handle(`sys/kbd/<down>`, func(e ui.Event) {
		log.TextFgColor = ui.ColorRed
	})

}
func handleActivity(ms *ui.Par, url string) {
	ms.Handle(`/timer/1000ms`, func(e ui.Event) {
		if isOnline(url) {
			ms.TextFgColor = ui.ColorGreen
		} else {
			ms.TextFgColor = ui.ColorRed
		}
	})
}

func service() {
	log.Width = 50
	log.Height = 10
	log.SetX(30)
	scripting.Width = 20
	scripting.Height = 3
	handleActivity(scripting, scriptingURL)

	tasking.Width = 20
	tasking.Height = 3
	tasking.SetY(3)
	handleActivity(tasking, taskingURL)

	agent.Width = 20
	agent.Height = 3
	agent.SetY(6)
	handleActivity(agent, agentURL)

}

func isOnline(url string) bool {
	req, err := http.NewRequest(`GET`, url, nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
