package strutil

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/parnurzeal/gorequest"
)

const (
	BondUri    = "http://180.169.168.210:8382"
	BondAppKey = "25df45508859405eb2444fc502e2b791"
	BondSecret = "48db3717ca6c44b99a6396544ebc7f92"
)

var (
	BondSign = base64.StdEncoding.EncodeToString([]byte(BondAppKey + " " + BondSecret))
)

func TestGetMaxHtml(t *testing.T) {

}

func TestGetBonds(t *testing.T) {
	fmt.Println(BondSign)
	uri := BondUri + "/ngpc/v1/bond/getBonds"
	a, b, c := gorequest.New().Get(uri).Set("Authorization", BondSign).End()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}

func TestGetVoteStats(t *testing.T) {
	fmt.Println(BondSign)
	uri := BondUri + "/ngpc/v1/predictionVote/getVoteStats?bondId=170210X5"
	a, b, c := gorequest.New().Get(uri).Set("Authorization", BondSign).End()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}

func TestGetLoadUploaded(t *testing.T) {
	uri := BondUri + "/ngpc/v1/bidResult/loadUploaded"
	a, b, c := gorequest.New().Get(uri).Set("Authorization", BondSign).End()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	//170204X9  170206X4  170210X5
}

func TestGetLoadMatched(t *testing.T) {
	uri := BondUri + "/ngpc/v1/bidResult/loadMatched"
	a, b, c := gorequest.New().Get(uri).Set("Authorization", BondSign).End()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	//170204X9  170206X4  170210X5
}

func TestRemoveEmptyPTag(t *testing.T) {
	s := "<p>每个人都可能成为一名“宽客”（QUANT）。</p><p><br> 《量化投资24小时》是我们推出的第一部“24小时”系投资训练特辑，包含视频及音频。希望能成为你第一部手把手“教程”。从零开始熟悉量化交易工具、经典策略搭建和更新迭代。带领个人投资者打开一个更具知识体系化的宽客世界。</p><p><br> <strong>【《量化投资24小时》的特训有什么独特之处】</strong></p><p>&nbsp;</p><p>第一、从最基础的python语言中就渗透入了交易的场景，比如列表在量化交易中到底是有什么起什么作用、字典又是扮演的额什么角色。场景学习，“边学边做”是最快速学习的方法。</p><p><br> 第二、主讲人，作为一个非计算机也非理工科背景出身“无编程经验”的金融硕士生，怎么能慢慢能成为一个宽客，并开发出vn.py这样一个国内用户最多的量化金融开源项目？ 这条路上会遭遇的挫折和你有着感同身受， 他说“希望把我自己所学分享给大家，也帮助更多的人少走一些我在早期时期走过的弯路。”</p><p>&nbsp;</p><p><strong>上周已更新视频</strong><br> 【11】循环语句：开启实时交易的无限循环<br> 【12】使用函数回测你的交易盈利<br> 【13】类和实例——组成大交易策略的零部件</p><p><br> <strong>本周主讲人准备线下见面活动，暂停更新一周。线下活动详情请见专辑最新文章</strong></p><p><br> 订阅读者可进入专属交流群，还将另行安排线上线下教学交流活动。</p>"
	fmt.Println(RemoveEmptyPTag(s))
}

func TestRemoveStylesOfHtmlTag(t *testing.T) {
	const sample_html  = `
 <!-- wp:paragraph -->
 <p><s>这是</s><strong class="tss-bold">一段</strong><em class="tss-italic">测试</em><span style="text-decoration: underline" class="tss-underline">文本</span></p>
 <!-- /wp:paragraph -->
 
 <!-- wp:paragraph -->
 <p><span style="background-color:#0059FF;font-style:#0059F;" class="tss-background-color"><span style="font-weight:#E60000" class="tss-color">这是一段测试文本</span></span></p>
 <!-- /wp:paragraph -->`

	fmt.Println(RemoveStylesOfHtmlTag(sample_html, "text-decoration", "font-weight", "font-style"))
}
