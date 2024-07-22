package global

type LanguageType uint

const (
	all LanguageType = iota + 1
	goLang
)

type Language string

const (
	All    Language = "all"
	GoLang Language = "go"
	Java   Language = "java"
	Python Language = "python"
	JS     Language = "javascript"
	TS     Language = "typescript"
	Ruby   Language = "ruby"
	PHP    Language = "php"
	Swift  Language = "swift"
	Kotlin Language = "kotlin"
	Rust   Language = "rust"
	Scala  Language = "scala"
)

var isPreview bool

func SetPreview(p bool) {
	isPreview = p
}
func IsPreviewMode() bool {
	return isPreview
}
func (l Language) LanguageType() LanguageType {
	switch l {
	case All:
		return all
	case GoLang:
		return goLang
	default:
		return all
	}
}
func (t LanguageType) Language() Language {
	switch t {
	case all:
		return "All"
	case goLang:
		return "Go"
	default:
		return "All"
	}
}
