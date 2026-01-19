package models

// modelRegistry 儲存所有需要 AutoMigrate 的模型
var modelRegistry []interface{}

// RegisterModel 在 init 中呼叫以註冊模型
func RegisterModel(model interface{}) {
	modelRegistry = append(modelRegistry, model)
}

// AllModels 回傳所有已註冊的模型
func AllModels() []interface{} {
	return modelRegistry
}
