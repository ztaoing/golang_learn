/**
* @Author:zhoutao
* @Date:2020/12/12 上午9:32
* @Desc:
 */

package bridging

func ExampleCommonSMS() {
	m := NewCommonMessage(ViaSMS())
	m.SendMessage("sending by sms:have logic drink", "boo")
}

func ExampleCommonEmail() {
	m := NewCommonMessage(ViaEmail())
	m.SendMessage("sending by email:have logic drink", "boo")
}

func ExampleUrgencySMS() {
	m := NewUrgencyMessage(ViaSMS())
	m.SendMessage("have logic drink", "boo")
}

func ExampleUrgencyEmail() {
	m := NewUrgencyMessage(ViaEmail())
	m.SendMessage("have logic drink", "boo")
}
