package utils

import (
	"fmt"
	"time"
)

type Formatter struct {
	owner string
}

func NewFormatter(owner string) *Formatter {
	return &Formatter{owner: owner}
}

func (f *Formatter) Message(msg string, owner string, sender string) string {
	// [geeks]  2023/03/12 14:09:28 ::: [dellno] hey dell
	formatted := fmt.Sprintf("[%v] %v ::: %v %v",
		owner,
		time.Now().Format("2006-01-02 15:04:05"),
		sender,
		msg)
	return formatted
}

func (f *Formatter) MessageCLI(msg string, owner string, sender string) string {
	// [geeks]  2023/03/12 14:09:28 ::: [dellno] hey dell
	formatted := fmt.Sprintf("\n[%v] %v ::: [%v] %v\n\n",
		owner,
		time.Now().Format("2006-01-02 15:04:05"),
		sender,
		msg)
	return formatted
}

func (f *Formatter) MessagePrimaryCLI(msg string, owner string, sender string) string {
	// [geeks]  2023/03/12 14:09:28 ::: [dellno] hey dell
	formatted := fmt.Sprintf("\n[%v] %v ::: %v %v\n",
		owner,
		time.Now().Format("2006-01-02 15:04:05"),
		sender,
		msg)
	return formatted
}

func (f *Formatter) MessageSecondaryCLI(msg string, owner string, sender string) string {
	// [geeks]  ................. ::: [dellno] hey dell
	formatted := fmt.Sprintf("[%v] %v ::: %v %v\n",
		owner,
		"...................",
		sender,
		msg)
	return formatted
}

func (f *Formatter) Log(msg string, owner string) string {
	// [Server] 2023/03/12 13:47:59 INFO ::: Server up and running on address: 0.0.0.0:3000
	formatted := fmt.Sprintf("[%v] %v %v ::: %v",
		f.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		owner,
		msg)
	return formatted
}
