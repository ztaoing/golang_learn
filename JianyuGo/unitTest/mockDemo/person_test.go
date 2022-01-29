/**
* @Author:zhoutao
* @Date:2022/1/29 14:38
* @Desc:
 */

package main

import (
	"golang_learn/golang_learn/JianyuGo/unitTest/mockDemo/equipment"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPerson_dayLife(t *testing.T) {
	type fields struct {
		name  string
		phone equipment.Phone
	}

	//生成mockphone对象
	mockCtl := gomock.NewController(t)
	mockPhone := equipment.NewMockPhone(mockCtl)
	//设置mockphone对象的接口方法返回参数
	mockPhone.EXPECT().ZhiHu().Return(true)
	mockPhone.EXPECT().WeiXin().Return(true)
	mockPhone.EXPECT().WangZhe().Return(true)

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"case1", fields{"iphone6s", equipment.NewIphone6s()}, true},
		{"case2", fields{"mocked phone", mockPhone}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Person{
				name:  tt.fields.name,
				phone: tt.fields.phone,
			}
			if got := x.dayLife(); got != tt.want {
				t.Errorf("dayLife() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 对接口进行mock，可以让我们在未实现具体对象的接口功能前，或者该接口调用代价非常高时，也能对业务代码进行测试。
// 而且在开发过程中，我们同样可以利用mock对象，不用因为等待接口实现方，实现相关功能，从而停止后续开发工作。
