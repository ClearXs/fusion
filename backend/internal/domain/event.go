package domain

type EventKey = string

type EventItem struct {
	EventName        EventKey `json:"EventName"`
	EventNameChinese string   `json:"EventNameChinese"`
	EventDescription string   `json:"EventDescription"`
	Passive          bool     `json:"Passive"`
}

const (
	LoginEvent               EventKey = "login"
	LogoutEvent              EventKey = "logout"
	BeforeUpdateArticleEvent EventKey = "beforeUpdateArticle"
	AfterUpdateArticleEvent  EventKey = "afterUpdateArticle"
	DeleteArticleEvent       EventKey = "deleteArticle"
	BeforeUpdateDraftEvent   EventKey = "beforeUpdateDraft"
	AfterUpdateDraftEvent    EventKey = "afterUpdateDraft"
	DeleteDraftEvent         EventKey = "deleteDraft"
	UpdateSiteInfoEvent      EventKey = "updateSiteInfo"
	ManualTriggerEvent       EventKey = "manualTriggerEvent"
)

var SystemEvents = []*EventItem{
	{
		EventName:        LoginEvent,
		EventNameChinese: "登录",
		EventDescription: "登录",
		Passive:          true,
	},
	{
		EventName:        LogoutEvent,
		EventDescription: "登出",
		EventNameChinese: "登出",
		Passive:          true,
	},
	{
		EventName:        BeforeUpdateArticleEvent,
		EventNameChinese: "更新文章之前",
		EventDescription: "更新文章之前，具体涉及到：发布草稿、保存文章、创建文章、更新文章信息，在此修改文章数据并返回会改变实际保存到数据库的值",
		Passive:          false,
	},
	{
		EventName:        AfterUpdateArticleEvent,
		EventNameChinese: "更新文章之后",
		EventDescription: "更新文章之后，具体涉及到：发布草稿、保存文章、创建文章、更新文章信息",
		Passive:          true,
	},
	{
		EventName:        DeleteArticleEvent,
		EventNameChinese: "删除文章",
		EventDescription: "删除文章",
		Passive:          true,
	},
	{
		EventName:        BeforeUpdateDraftEvent,
		EventNameChinese: "更新草稿之前",
		EventDescription: "更新草稿之前，具体涉及到：保存草稿、创建草稿、更新草稿信息，在此修改文章内容并返回会改变实际保存到数据库的文章内容",
		Passive:          false,
	},
	{
		EventName:        AfterUpdateDraftEvent,
		EventNameChinese: "更新草稿之后",
		EventDescription: "更新草稿之后，具体涉及到：保存草稿、创建草稿、更新草稿信息",
		Passive:          true,
	},
	{
		EventName:        DeleteDraftEvent,
		EventNameChinese: "删除草稿",
		EventDescription: "删除草稿",
		Passive:          true,
	},
	{
		EventName:        UpdateSiteInfoEvent,
		EventNameChinese: "更新站点信息",
		EventDescription: "更新站点信息",
		Passive:          true,
	},
	{
		EventName:        ManualTriggerEvent,
		EventNameChinese: "手动触发事件",
		EventDescription: "手动触发事件事件",
		Passive:          true,
	},
}
