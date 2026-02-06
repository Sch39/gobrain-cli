package config

type Config struct {
	Version    string              `mapstructure:"version"`
	Debug      bool                `mapstructure:"debug"`
	Project    Project             `mapstructure:"project"`
	Meta       Meta                `mapstructure:"meta"`
	Tools      []Tool              `mapstructure:"tools"`
	Scripts    map[string][]string `mapstructure:"scripts"`
	Generators Generators          `mapstructure:"generators"`
	Verify     verify              `mapstructure:"verify"`
	Layout     Layout              `mapstructure:"layout"`
}

type Project struct {
	Name      string `mapstructure:"name"`
	Module    string `mapstructure:"module"`
	Toolchain string `mapstructure:"toolchain"`
}

type Meta struct {
	TemplateOrigin string `mapstructure:"template_origin"`
	Author         string `mapstructure:"author"`
}

type Tool struct {
	Name string `mapstructure:"name"`
	Pkg  string `mapstructure:"pkg"`
}

type Generators struct {
	Requires []string                    `mapstructure:"requires"`
	Commands map[string]GeneratorCommand `mapstructure:"commands"`
}

type GeneratorCommand struct {
	Desc      string     `mapstructure:"desc"`
	Args      []string   `mapstructure:"args"`
	Templates []Template `mapstructure:"templates"`
}

type Template struct {
	Src  string `mapstructure:"src"`
	Dest string `mapstructure:"dest"`
}

type verify struct {
	FailFast bool         `mapstructure:"fail_fast"`
	Pipeline []VerifyStep `mapstructure:"pipeline"`
}

type VerifyStep struct {
	Name string `mapstructure:"name"`
	Run  string `mapstructure:"run"`
}

type Layout struct {
	Dirs []string `mapstructure:"dirs"`
}
