package main

type MockRule struct {
	MockEntityAdapter interface{}
	OtherAdapter      interface{}
	OtherRule         interface{}
}

type MockEntityRepo struct {
}

type MockUseCase struct {
	MockRuler    interface{}
	DBAdapter    interface{}
	QueryAdapter interface{}
}

func NewMockServiceServer() {

}
