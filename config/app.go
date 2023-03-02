package config

type Application struct {
	Env *Env
}

func App(env string) Application {
	app := &Application{}
	app.Env = NewEnv(env)
	return *app
}
