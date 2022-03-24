package ticket

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const cronTime = "*/3 * * * * ?"
const ticketId = "955317580770140164467206"

func GetTicketNumber() {
	ch := make(chan bool)
	c := cron.New()
	err := c.AddFunc(cronTime, func() {
		r := rand.Intn(3000)
		time.Sleep(time.Duration(r) * time.Millisecond)
		b := GetData("https://kylin.zhengzai.tv/kylin/performance/" + ticketId)
		if b {
			ch <- true
		}
	})
	if err != nil {
		log.Printf("定时任务启动失败")
		log.Fatalln(err.Error())
		return
	}
	c.Start()
	select {
	case r := <-ch:
		if r {
			log.Println("有余票了")
		}
	}
	//PostMethod()
}

//这是get请求
func GetData(ulr string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //如果需要测试自签名的证书 这里需要设置跳过证书检测 否则编译报错
	client := &http.Client{Transport: tr}
	//提交请求
	reqest, err := http.NewRequest("GET", ulr, nil)
	//增加header选项
	reqest.Header.Add("Host", "kylin.zhengzai.tv")
	reqest.Header.Add("Origin", "https://m.zhengzai.tv")
	reqest.Header.Add("Referer", "https://m.zhengzai.tv/")
	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	defer response.Body.Close()
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	m := make(map[string]interface{})
	if err := json.Unmarshal(body, &m); err != nil {
		fmt.Println("error:", err)
		return false
	}
	data := m["data"].(map[string]interface{})
	ttL := data["ticketTimeList"].([]interface{})
	log.Printf("=============%v:开始进行判断===========================》", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range ttL {
		tl := v.(map[string]interface{})["ticketList"].([]interface{})

		for _, i := range tl {
			j := i.(map[string]interface{})
			log.Printf("开始判断时间:%v，价格:%v的余票", j["useStart"], j["discountPrice"])
			c := (int)(j["status"].(float64))
			if c != 8 && c != 10 {
				fmt.Printf("有余票了：时间：%v，价格：%v \n", j["useStart"], j["discountPrice"])
				return true
			}
		}
	}
	log.Printf("没有余票")
	log.Printf("==============%v:判断结束=======================》", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("")
	return false
}
