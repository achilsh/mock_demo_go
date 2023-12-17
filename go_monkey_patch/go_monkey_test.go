package gomonkeypatch

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)


func CallFunc(ctx context.Context, a int, b int)  (int, error) {
	if a < 0 || b < 0 {
		return 0, errors.New("is invalid param")
	}
	return a+b, nil
}

type DemoGomonkeyMethod struct {

}

func (d *DemoGomonkeyMethod) CallA(ctx context.Context, a int, b int) (c int, e error) {
	return a+b, nil
}

func TestGoMonkeyDemoForFunc(t *testing.T) {
	
	{
		p := gomonkey.NewPatches()
		//第一个参数是 真实函数名， 第二参数是 桩函数，
		p.ApplyFunc(CallFunc, func( context.Context, int, int) (int, error) {
			return 11+12, nil
		})
		defer p.Reset()

		//调用真实函数时，其实调用桩函数。
		a, _ := CallFunc(context.Background(), 1,2)
		t.Logf("callFunc(), ret: %v", a)
	}
	
	{
		//类型的 method打桩， 用于模拟struct 的成员方法。
		p := gomonkey.NewPatches()
		b := &DemoGomonkeyMethod{}
		
		var c *DemoGomonkeyMethod
		//第一个参数是 对象的类型，一般用反射来获取，第二参数 是menthod 字符串名， 第三个参数就是桩 方法，其中桩方法的第一个参数是：struct的receiver，其他是真实函数的参数和返回值，
		p.ApplyMethod(reflect.TypeOf(c), "CallA", func( _ *DemoGomonkeyMethod,_ context.Context, _ int, _ int)(int, error){
			return 100+200, nil
		}) 
		defer p.Reset()
		
		//调用真实方法，其实调用桩方法，比如： 
		v, _ := b.CallA(context.Background(),0,0)
		t.Logf("call method() ret: %v", v)
	}


}