/**
* @Author:zhoutao
* @Date:2023/1/3 16:20
* @Desc: 抽象工厂 用于创建一系列相关的或者相互依赖的对象。
 */

package main

/**
目前抽象工厂有两个实际工厂类一个是华为的工厂，一个是小米的工厂，他们用来实际生产自家的产品设备。
*/

func main() {
	/**
	抽象工厂模式与工厂方法模式最大的区别在于，工厂方法模式针对的是一个产品等级结构，
	而抽象工厂模式则需要面对多个产品等级结构，一个工厂等级结构可以负责多个不同产品等级结构中的产品对象的创建 。
	*/
	var factory AbstractFactory
	var tv ITelevision
	var air IAirConditioner

	factory = &HuaWeiFactory{}
	tv = factory.CreateTelevision()
	air = factory.CreateAirConditioner()
	tv.Watch()
	air.SetTemperature(25)

	factory = &MiFactory{}
	tv = factory.CreateTelevision()
	air = factory.CreateAirConditioner()
	tv.Watch()
	air.SetTemperature(26)
}
