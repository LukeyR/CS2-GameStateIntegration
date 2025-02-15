package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

func main() {
	http.HandleFunc("/bombtimer", bombTimerHandler)
	log.Println("Server listening on :8000")

	cs2gsi.RegisterEventHandler(events.EventBombPlanted, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Println("The Bomb has been planted", gsiEvent.GetOriginalRequestFlat())
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponRemoved, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Println("The Bomb dropped", gsiEvent.GetOriginalRequestFlat())
	})

	cs2gsi.RegisterGlobalHandler(func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		fmt.Printf("(%v) %v %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			events.EnumToEventName[gameEvent.EventType],
			gsiEvent.GetOriginalRequestFlat(),
		)
	})

	err := cs2gsi.StartupAndServe(":8000")
	if err != nil {
		return
	}
}

//go:embed index.html
var templates embed.FS

func bombTimerHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(&templates, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
