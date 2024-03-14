package domain

import (
	"time"
)

type SocialType = string
type TrueOrFalse = string
type Theme = string
type Site = string

const (
	TrueString  TrueOrFalse = "true"
	FalseString TrueOrFalse = "false"
)

const (
	AutoTheme  Theme = "auto"
	DarkTheme  Theme = "dark"
	LightTheme Theme = "light"
)

const (
	LogoSite Site = "siteLogo"
	NameSite Site = "siteName"
)

const (
	BilibiliSocialType   SocialType = "bilibili"
	EmailSocialType      SocialType = "email"
	GithubSocialType     SocialType = "github"
	GiteeSocialType      SocialType = "gitee"
	WechatSocialType     SocialType = "wechat"
	WechatDarkSocialType SocialType = "wechat-dark"
)

var DefaultSocials = []*SocialItem{
	{
		Value: "哔哩哔哩",
		Type:  BilibiliSocialType,
	},
	{
		Value: "邮箱",
		Type:  EmailSocialType,
	},
	{
		Value: "GitHub",
		Type:  GithubSocialType,
	},
	{
		Value: "Gitee",
		Type:  GiteeSocialType,
	},
	{
		Value: "微信",
		Type:  WechatSocialType,
	},
	{
		Value: "微信（暗色模式）",
		Type:  WechatDarkSocialType,
	},
}

var DefaultMenu = []MenuItem{
	{
		Id:    0,
		Name:  "首页",
		Value: "/",
		Level: 0,
	},
	{
		Id:    1,
		Name:  "标签",
		Value: "/tag",
		Level: 0,
	},
	{
		Id:    2,
		Name:  "分类",
		Value: "/category",
		Level: 0,
	},
	{
		Id:    3,
		Name:  "时间线",
		Value: "/timeline",
		Level: 0,
	},
	{
		Id:    4,
		Name:  "友链",
		Value: "/link",
		Level: 0,
	},
	{
		Id:    5,
		Name:  "关于",
		Value: "/about",
		Level: 0,
	}}

type Meta struct {
	Id             int64         `json:"id" bson:"id"`
	Links          []*LinkItem   `json:"links" bson:"links"`
	Socials        []*SocialItem `json:"socials" bson:"socials"`
	Menus          []*MenuItem   `json:"menus" bson:"menus"`
	Rewards        []*RewardItem `json:"rewards" bson:"rewards"`
	About          *About        `json:"about" bson:"about"`
	SiteInfo       *SiteInfo     `json:"siteInfo" bson:"siteInfo"`
	Viewer         int64         `json:"viewer" bson:"viewer"`
	Visited        int64         `json:"visited" bson:"visited"`
	Categories     []string      `json:"categories" bson:"categories"`
	TotalWordCount int64         `json:"totalWordCount" bson:"totalWordCount"`
}

type LinkItem struct {
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Url       string    `json:"url" bson:"url"`
	Name      string    `json:"name" bson:"name"`
	Desc      string    `json:"desc" bson:"desc"`
	Logo      string    `json:"logo" bson:"logo"`
}

type SocialItem struct {
	UpdatedAt time.Time  `json:"updatedAt" bson:"updatedAt"`
	Value     string     `json:"value" bson:"value"`
	Type      SocialType `json:"type" bson:"type"`
}

type MenuItem struct {
	Id       int64     `json:"id" bson:"id"`
	Name     string    `json:"name" bson:"name"`
	Value    string    `json:"value" bson:"value"`
	Level    int32     `json:"level" bson:"level"`
	Children *MenuItem `json:"children" bson:"children"`
}

type RewardItem struct {
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Value     string    `json:"value" bson:"value"`
	Name      string    `json:"name" bson:"name"`
}

type About struct {
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Content   string    `json:"content" bson:"content"`
}

type SiteInfo struct {
	Author                      string      `json:"author" bson:"author"`
	AuthorLogo                  string      `json:"authorLogo" bson:"authorLogo"`
	AuthorLogoDark              string      `json:"authorLogoDark" bson:"authorLogoDark"`
	AuthDesc                    string      `json:"authDesc" bson:"authDesc"`
	SiteLogo                    string      `json:"siteLogo" bson:"siteLogo"`
	SiteLogoDark                string      `json:"siteLogoDark" bson:"siteLogoDark"`
	Favicon                     string      `json:"favicon" bson:"favicon"`
	SiteName                    string      `json:"siteName" bson:"siteName"`
	SiteDesc                    string      `json:"siteDesc" bson:"siteDesc"`
	BeianNumber                 string      `json:"beianNumber" bson:"beianNumber"`
	BeianUrl                    string      `json:"beianUrl" bson:"beianUrl"`
	GaBeianNumber               string      `json:"gaBeianNumber" bson:"gaBeianNumber"`
	GaBeianUrl                  string      `json:"gaBeianUrl" bson:"gaBeianUrl"`
	GaBeianLogoUrl              string      `json:"gaBeianLogoUrl" bson:"gaBeianLogoUrl"`
	PayAliPay                   string      `json:"payAliPay" bson:"payAliPay"`
	PayWechat                   string      `json:"payWechat" bson:"payWechat"`
	PayAliPayDark               string      `json:"payAliPayDark" bson:"payAliPayDark"`
	PayWechatDark               string      `json:"payWechatDark" bson:"payWechatDark"`
	Since                       time.Time   `json:"since" bson:"since"`
	BaseUrl                     string      `json:"baseUrl" bson:"baseUrl"`
	GaAnalysisId                string      `json:"gaAnalysisId" bson:"gaAnalysisId"`
	BaiduAnalysisId             string      `json:"baiduAnalysisId" bson:"baiduAnalysisId"`
	CopyrightAggreement         string      `json:"copyrightAggreement" bson:"copyrightAggreement"`
	EnableComment               TrueOrFalse `json:"enableComment" bson:"enableComment"`
	ShowSubMenu                 TrueOrFalse `json:"showSubMenu" bson:"showSubMenu"`
	HeaderLeftContent           Site        `json:"headerLeftContent" bson:"headerLeftContent"`
	SubMenuOffset               int64       `json:"subMenuOffset" bson:"subMenuOffset"`
	ShowAdminButton             TrueOrFalse `json:"showAdminButton" bson:"showAdminButton"`
	ShowDonateInfo              TrueOrFalse `json:"showDonateInfo" bson:"showDonateInfo"`
	ShowFriends                 TrueOrFalse `json:"showFriends" bson:"showFriends"`
	ShowCopyRight               TrueOrFalse `json:"showCopyRight" bson:"showCopyRight"`
	ShowDonateButton            TrueOrFalse `json:"showDonateButton" bson:"showDonateButton"`
	ShowDonateInAbout           TrueOrFalse `json:"showDonateInAbout" bson:"showDonateInAbout"`
	AllowOpenHiddenPostByUrl    TrueOrFalse `json:"allowOpenHiddenPostByUrl" bson:"allowOpenHiddenPostByUrl"`
	DefaultTheme                Theme       `json:"defaultTheme" bson:"defaultTheme"`
	EnableCustomizing           TrueOrFalse `json:"enableCustomizing" bson:"enableCustomizing"`
	ShowRSS                     TrueOrFalse `json:"showRSS" bson:"showRSS"`
	OpenArticleLinksInNewWindow TrueOrFalse `json:"openArticleLinksInNewWindow" bson:"openArticleLinksInNewWindow"`
	ShowExpirationReminder      TrueOrFalse `json:"showExpirationReminder" bson:"showExpirationReminder"`
	ShowEditButton              TrueOrFalse `json:"showEditButton" bson:"showEditButton"`
}
