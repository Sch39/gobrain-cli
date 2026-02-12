package config

type Config struct {
	Version    string              `mapstructure:"version" yaml:"version"`
	Debug      bool                `mapstructure:"debug" yaml:"debug"`
	Project    Project             `mapstructure:"project" yaml:"project"`
	Meta       Meta                `mapstructure:"meta" yaml:"meta"`
	Tools      []Tool              `mapstructure:"tools" yaml:"tools"`
	Scripts    map[string][]string `mapstructure:"scripts" yaml:"scripts"`
	Generators Generators          `mapstructure:"generators" yaml:"generators"`
	Verify     verify              `mapstructure:"verify" yaml:"verify"`
	Layout     Layout              `mapstructure:"layout" yaml:"layout"`
}

type Project struct {
	Name      string `mapstructure:"name" yaml:"name"`
	Module    string `mapstructure:"module" yaml:"module"`
	Toolchain string `mapstructure:"toolchain" yaml:"toolchain"`
}

type Meta struct {
	TemplateOrigin string `mapstructure:"template_origin" yaml:"template_origin"`
	Author         string `mapstructure:"author" yaml:"author"`
}

type Tool struct {
	Name string `mapstructure:"name" yaml:"name"`
	Pkg  string `mapstructure:"pkg" yaml:"pkg"`
}

type Generators struct {
	Requires []string                    `mapstructure:"requires" yaml:"requires"`
	Commands map[string]GeneratorCommand `mapstructure:"commands" yaml:"commands"`
}

type GeneratorCommand struct {
	Desc      string     `mapstructure:"desc" yaml:"desc"`
	Args      []string   `mapstructure:"args" yaml:"args"`
	Templates []Template `mapstructure:"templates" yaml:"templates"`
}

type Template struct {
	Src  string `mapstructure:"src" yaml:"src"`
	Dest string `mapstructure:"dest" yaml:"dest"`
}

type verify struct {
	FailFast bool         `mapstructure:"fail_fast" yaml:"fail_fast"`
	Pipeline []VerifyStep `mapstructure:"pipeline" yaml:"pipeline"`
}

type VerifyStep struct {
	Name string `mapstructure:"name" yaml:"name"`
	Run  string `mapstructure:"run" yaml:"run"`
}

type Layout struct {
	Dirs []string `mapstructure:"dirs" yaml:"dirs"`
}
