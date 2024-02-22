package main

import "github.com/hyperf/gotask/v2/pkg/gotask"

type App struct{}

func (a *App) Hi(name string, r *interface{}) error {
	*r = map[string]string{
		"hello": name,
	}
	return nil
}

func main() {
	gotask.SetAddress("127.0.0.1:6001")
	gotask.Register(new(App))
	gotask.Run()
}
