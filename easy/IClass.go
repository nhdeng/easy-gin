package easy

type IClass interface {
	Build(easy *Easy)
	Name() string
}
