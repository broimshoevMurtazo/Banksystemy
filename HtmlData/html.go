package htmldata

import (
	"fmt"
	"html/template"
	"nank/app/structs"
	"os"
)

// Генерация HTML для email
func HtmlData(data structs.Registration) error {
	// Парсим шаблон
	tmpl, err := template.ParseFiles("./controlers/htmls/mail.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err) // Логируем ошибку при парсинге шаблона
		return err
	}

	// Путь к выходному файлу
	outputPath := "output.html"

	// Открываем файл для записи
	outputFile, err := os.Create("controlers" + "/" + "htmls"+"/"+outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err) // Логируем ошибку при создании файла
		return err
	}
	defer outputFile.Close()

	// Выполняем шаблон с данными
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err) // Логируем ошибку при выполнении шаблона
		return err
	}

	fmt.Println("HTML file generated successfully at:", outputPath) // Логируем успешное создание HTML файла
	return nil
}




func HtmlData2(Data structs.Registration) error {
	// Парсим шаблон
	tmpl, err := template.ParseFiles("./controlers/htmls/password.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err) // Логируем ошибку при парсинге шаблона
		return err
	}

	// Путь к выходному файлу
	outputPath := "passwordoutput.html"

	// Открываем файл для записи
	outputFile, err := os.Create("controlers" + "/" + "htmls"+"/"+outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err) // Логируем ошибку при создании файла
		return err
	}
	defer outputFile.Close()

	// Выполняем шаблон с данными
	err = tmpl.Execute(outputFile, Data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err) // Логируем ошибку при выполнении шаблона
		return err
	}

	fmt.Println("HTML file generated successfully at:", outputPath) // Логируем успешное создание HTML файла
	return nil
}