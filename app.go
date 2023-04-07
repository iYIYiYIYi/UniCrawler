package main

import (
	crawler "UniCrawler/cmd/crawlerr"
	"UniCrawler/cmd/util"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	util.InitDatabase()
	crawler.DefaultInit()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetAllSchools() []util.School {
	return *util.ReadAllSchools()
}

func (a *App) GetAfterTime(time string) []util.School {
	t,e := util.StrToTime(time)
	if e != nil {
		return nil
	}
	return *util.ReadAfterTime(&t)
}

func (a *App) GetWithKeyWords(keywords []string) []util.School {
	return *util.ReadSchoolsWithKeywords(keywords)
}

func (a *App) StartCrawler() {
	crawler.Start()
}

