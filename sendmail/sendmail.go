package sendmail

import (
	// "bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

func Sendmail() {
	// Генерация HTML-данных
	err := HtmlData()
	if err != nil {
		panic(err)
	}

	// Отправка email
	// err = SendGomail("output.html", "Hello from Go!")
	// if err != nil {
	// 	panic(err)
	// }
}
func randomSixDigit() int {
	// Устанавливаем seed для случайного генератора
	rand.Seed(time.Now().UnixNano())
	// Генерируем случайное число от 100000 до 999999
	return rand.Intn(900000) + 100000
}

func HtmlData() error {
	// Парсим HTML-шаблон
	tmpl := template.Must(template.ParseFiles("mail.html"))
    code:=randomSixDigit()
	// Данные для вставки в шаблон
	data := struct {
		Code int
	}{
		Code: code,
	}

	// Открываем файл для записи
	outputFile, err := os.Create("output.html")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Выполняем шаблон с передачей данных
	err = tmpl.ExecuteTemplate(outputFile, "mail.html", data)
	if err != nil {
		return err
	}

	return nil
}

func SendGomail(templatePath string, subject string , email  string) error {
	htmlData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "murtazobroimshoevm4@gmail.com")
	m.SetHeader("To",email) 
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", string(htmlData))

	d := gomail.NewDialer("smtp.gmail.com", 587, "murtazobroimshoevm4@gmail.com", "odlx sdkt aamz bjjn")

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	fmt.Println("Email sent successfully!")

	return nil

}
