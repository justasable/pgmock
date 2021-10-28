package generate

var Test_newGenerator = newGenerator

type Test_defaultValType = defaultValType

func (g generator) Test_TestVals() []interface{} {
	return g.testVals
}

func (g generator) Test_UniqueVal(n int) interface{} {
	return g.uniqueFn(n)
}
