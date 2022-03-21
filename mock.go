package main

type MockRule struct {
	MockEntityAdapter interface{}
	OtherAdapter      interface{}
}

type MockEntityRepo struct {
}

type MockUseCase struct {
	MockRuler    interface{}
	DBAdapter    interface{}
	QueryAdapter interface{}
}
