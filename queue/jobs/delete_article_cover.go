package jobs

type DeleteArticleCoverArgs struct {
	ArticleID int32  `json:"article_id"`
	Key       string `json:"key"`
}

func (DeleteArticleCoverArgs) Kind() string { return "delete_article_cover" }
