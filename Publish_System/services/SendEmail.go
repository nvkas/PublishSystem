package services

import (
	"github.com/go-gomail/gomail"
	"log"
)

func SendEmail(message ,formAddress,receiveAddress ,username,authorizationCode string)  {
	// 创建新的消息对象
	msg := gomail.NewMessage()

	//发件人
	msg.SetAddressHeader("From",formAddress,"错误报警")

	//收件人
	msg.SetHeader("To",msg.FormatAddress(receiveAddress,""))

	//正文
	msg.SetBody("text/html",message)

	//附件 msg.Attach

	//默认是QQ邮箱,需要进入邮箱“设置”,在“帐户”项里找到“POP3/SMTP服务”的设置项，进行开启
	d := gomail.NewDialer("smtp.qq.com", 465, username, authorizationCode) // 发送邮件服务器、端口、发件人账号、授权码
	if err := d.DialAndSend(msg); err != nil {
		log.Println("发送失败", err)
		return
	}

	return

}
