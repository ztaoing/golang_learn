/**
* @Author:zhoutao
* @Date:2023/1/3 18:38
* @Desc:
 */

package main

// 抽象工厂
type AbstractFactory interface {
	CreateTelevision() ITelevision
	CreateAirConditioner() IAirConditioner
}

type ITelevision interface {
	Watch()
}

type IAirConditioner interface {
	SetTemperature(int)
}
