package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

var content = `
curl 'https://www.baidu.com/sugrec?pre=1&p=3&ie=utf-8&json=1&prod=pc&from=pc_web&sugsid=39843,39934,39937,39932,39942,39938,39732,39991,39662,40009,40042&wd=x&his=%5B%7B%22time%22%3A1703491220%2C%22kw%22%3A%22proto%20%E7%BB%A7%E6%89%BF%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703491230%2C%22kw%22%3A%22protobuf%E7%BB%A7%E6%89%BF%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703492426%2C%22kw%22%3A%22%E7%BD%97%E7%A5%A5%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703492531%2C%22kw%22%3A%22golang%2021%20cmp%22%2C%22fq%22%3A3%7D%2C%7B%22time%22%3A1703492533%2C%22kw%22%3A%22%E8%AE%A2%E5%A9%9A%E5%BC%BA%E5%A5%B8%E6%A1%88%E4%B8%80%E5%AE%A1%E5%AE%A3%E5%88%A4%3A%E7%94%B7%E5%AD%90%E8%8E%B7%E5%88%913%E5%B9%B4%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703493248%2C%22kw%22%3A%22%E6%AD%A2%E5%92%B3%E5%AE%9D%E7%89%87%E6%9C%89%E7%BD%82%E7%B2%9F%E5%A3%B3%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703493411%2C%22kw%22%3A%22%E5%A9%9A%E5%89%8D%E5%8D%8F%E8%AE%AE%E6%9C%89%E6%B3%95%E5%BE%8B%E6%95%88%E5%8A%9B%E5%90%97%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703499683%2C%22kw%22%3A%22graph%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703499697%2C%22kw%22%3A%22graph%20http%22%2C%22fq%22%3A2%7D%2C%7B%22time%22%3A1703554479%2C%22kw%22%3A%22%E9%87%8E%E6%8C%87%E9%92%88%22%2C%22fq%22%3A2%7D%5D&req=2&csor=1&cb=jQuery110202410573158026359_1703557210985&_=1703557210986' \
  -H 'Accept: text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01' \
  -H 'Accept-Language: zh-CN,zh-TW;q=0.9,zh;q=0.8,en-US;q=0.7,en;q=0.6' \
  -H 'Connection: keep-alive' \
  -H 'Cookie: BIDUPSID=025110782CD25E32136CD4CF332B853D; PSTM=1687921565; BD_UPN=123253; BAIDUID=294BCEABD0D2E077FEE0BCDF31AFEB51:FG=1; H_WISE_SIDS_BFESS=131862_114550_216846_213356_214804_110085_244729_254835_259031_261719_236312_256419_265881_266360_265615_267072_266714_268566_268592_266186_269022_259642_256154_269730_269778_268235_269749_269904_269770_270083_267066_256739_270336_270460_270603_270547_270922_271023_271173_271178_271075_268987_269034_271228_267659_271320_271350_270279_270102_271562_270442_270157_271813_269875_271938_271958_269665_271949_271254_234295_269296_271188_272282_266565_267596_272321_272364_272009_272337_272466_272614_253022_272659_271688_272608_272821_272816_272801_260335_272986_269715_273061_267559_273093_273118_273131_273147_273243_273300_273400_273381_271158_270055_273525_273529_273521_273514_272641_273463_273197_272561_271147_273671_273705_264170_270186_263619_273164_274080_273960_273965_274141_269609_273917_274238_273788; sugstore=1; H_WISE_SIDS=39712_39843_39903_39819_39909_39934_39937_39932_39942_39940_39938_39930_39732_39662_39962; MCITY=-340%3A; H_PS_PSSID=39843_39934_39937_39932_39942_39938_39732_39991_39662_40009_40042; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BA_HECTOR=2k850h2h8gaga4208g2gah241iohrob1q; delPer=0; BD_CK_SAM=1; PSINO=6; ZFY=ahhfZeitHlk3soN6MjFMJWaDXJRGznDee3b:AXr:B8t1g:C; BAIDUID_BFESS=294BCEABD0D2E077FEE0BCDF31AFEB51:FG=1; baikeVisitId=e1e44239-056c-4ae4-8b5b-6d455240fab4; RT="z=1&dm=baidu.com&si=8d72000f-a8fe-4f94-82a1-4ffbd6ccda59&ss=lqknj6w8&sl=1&tt=1vk&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&nu=9y8m6cy&cl=2fg&ld=2nx&ul=cav&hd=cbd"; COOKIE_SESSION=262457_0_8_8_8_21_1_0_7_6_1_1_20794_0_0_0_1703230995_0_1703493394%7C9%230_1_1702372365%7C1' \
  -H 'Ps-Dataurlconfigqid: 0xa70d624300089bfa' \
  -H 'Referer: https://www.baidu.com/' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36' \
  -H 'X-Requested-With: XMLHttpRequest' \
  -H 'sec-ch-ua: "Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  --compressed
`

func TestParseCurl(t *testing.T) {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimSuffix(content, "\n")
	lines := strings.Split(content, "\n")

	regData, _ := regexp.Compile(`'(.*?)'`)
	findResult := regData.FindStringSubmatch(lines[0])
	_url := findResult[1]
	fmt.Println(_url)

	header := http.Header{}
	for _, s := range lines[1:] {
		data := regData.FindStringSubmatch(s)
		if len(data) != 2 {
			continue
		}
		part := strings.SplitN(data[1], ":", 2)
		key := part[0]
		value := strings.TrimSpace(part[1])
		header.Set(key, value)
	}

	fmt.Println(header)
}
