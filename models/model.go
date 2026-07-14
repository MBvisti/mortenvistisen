// Package models contains data models and validation logic.
package models

type (
	user                 struct{}
	token                struct{}
	article              struct{}
	tag                  struct{}
	articleTagConnection struct{}
	newsletter           struct{}
	subscriber           struct{}
	project              struct{}
	job                  struct{}
)

var (
	User                 user
	Token                token
	Article              article
	Tag                  tag
	ArticleTagConnection articleTagConnection
	Newsletter           newsletter
	Subscriber           subscriber
	Project              project
	Job                  job
)
