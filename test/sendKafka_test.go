package test

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"zhibo/kafka"
)

var text = "【趋势赛道梳理】\n1.赛道股\n1）新能源汽车。整车，长安汽车继续新高，其他，江淮汽车这种相对低位品种；锂矿和中游都一般\n2）光伏。上游核心股大跌，石英股份，高测股份都跌到过7%+，前几天上游补涨到格局较差的有机硅，金属硅，今天下跌如新安股份的大阴，需要开始留意上游的负反馈。 板块逆势：固德威\n3）半导体，之前一直是赛道最弱的，今天起来了。主要是汽车半导体，博通集成，富含维，国芯科技，闻泰科技等。\n2.周期股\n猪肉继续反弹，傲农生物，新五丰，唐人神\n煤炭，晋控煤业新高\n磷化工，整体调整，云天化高开低走但跌幅不算深\n3.大金融\n证券，长城证券，光大证券\n保险，中国人寿，中国人保，中国太保\n银行，瑞丰银行，张家港行，兰州银行，平安银行\n4.军工。中船系，中船防务，中船科技，中国船舶等\n5. 基建，山东路桥，信达地产等首板带队"

func TestSendKafka1(t *testing.T) {
	product := &kafka.Product{Config: kafka.Config{Address: "49.233.18.215:9092"}}
	product.Instance()
	product.Push("group-1",
		kafka.InitMessage(
			text,
			"",
			"anno",
			"734414412644999168",
			"2022-06-15 11:41:46").ToJson())
}

func TestMakeStr(t *testing.T) {
	ti, _ := time.Parse("2006-01-02T15:04:05.000+0800", "2022-11-04T14:07:15.435+0800")
	fmt.Println(ti.Format("2006-01-02 15:04"))
}

func makeStr(str []string) []string {
	var newStr []string
	for _, s := range str {
		if strings.Count(s, "") < 20 {
			newStr = append(newStr, s)
		} else {
			i := 0
			re := ""
			for _, ss := range strings.Split(s, "") {
				if i > 20 {
					newStr = append(newStr, re)
					i = 0
					re = ""
				}
				i++
				re += ss
			}
			newStr = append(newStr, re)
		}
	}
	return newStr
}
