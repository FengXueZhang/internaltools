package main

import (
	"fmt"
	"colly"
  "os"
  "strings"
	"mail"
  "io/ioutil"
)

// 追加写文件内容
func TraceFile(filename string, str_content string) {
    fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    fd_content := strings.Join([]string{str_content}, "")
    buf := []byte(fd_content)
    fd.Write(buf)
    fd.Close()
}

// 读取文件内容
func ReadFile(file_name string) (string, error) {
  b, err := ioutil.ReadFile(file_name)
  if err != nil {
      return "", err
  }
  str := string(b)
	return str, err
}

func main() {

	// 这里支持map类型，前提是目标地址页面结构一致
	myMap := map[string]string {
		"https://blogs.technet.microsoft.com/office_sustained_engineering/" : "CVE_Office",
		}

	c := colly.NewCollector()

	// 这里发起请求
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %v\n", r.URL)
	})

	// 这里是匹配html标签
  c.OnHTML("div.site-content", func(e *colly.HTMLElement) {

			key := e.Request.URL.Scheme + "://" + e.Request.URL.Host + e.Request.URL.Path
			file_name := myMap[string(key)] + "_temp.txt"

      TraceFile(file_name, e.Text)
  })

	// 这里是爬取结束的回调函数
  c.OnScraped(func(r *colly.Response) {
		key := r.Request.URL.Scheme + "://" + r.Request.URL.Host + r.Request.URL.Path
		base_file_name := myMap[string(key)] + ".txt"
		temp_file_name := myMap[string(key)] + "_temp.txt"

		base_file_text,err := ReadFile(base_file_name)
		if err != nil {
	      fmt.Println("OnScraped:", err, "\n")
	  }
		temp_file_text,err := ReadFile(temp_file_name)
		if err != nil {
	      fmt.Println("OnScraped:", err, "\n")
	  }

		// 对比字符串
		IsEqual := strings.EqualFold(base_file_text, temp_file_text)

		if IsEqual != true {
				fmt.Println("邮件告警\n")
				// 发邮件
				email := mail.NewEMail(`{"port":25}`)
				email.From = `*******@**.**`
				email.Host = `**.**.**.**`
				email.Port = 25
				email.Username = `****`
				email.Password = `****`
				email.Auth = mail.LoginAuth(email.Username, email.Password)

				email.To = []string{"mr.zhangxuefeng@gmail.com",}
				email.Subject = "安全小爬虫"
				email.Text = "尊敬的用户：\r\n 漏洞信息发生变化:" + key
				err := email.Send()
				if err != nil {
					fmt.Println(err)
				}
		}

		// 将临时文件替换为对比基准文件
		ioutil.WriteFile(base_file_name, []byte(temp_file_text), 0644)
		os.Remove(temp_file_name)
    fmt.Print("Finished\n")
  })

	// 循环需要爬取得map
	for i := range myMap {
		c.Visit(i)
	}

}
