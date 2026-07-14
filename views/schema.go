package views

import (
	"maps"
	"strings"
	"time"

	"mortenvistisen/models"
	"mortenvistisen/router/routes"
)

type SchemaNode map[string]any
type SchemaBuilder func(HeadData) []SchemaNode

type SchemaBreadcrumb struct {
	Name string
	Path string
}

type schemaListItem struct {
	Name        string
	Path        string
	Description string
	Image       string
}

var personSameAs = []string{
	"https://x.com/mbvisti",
	"https://github.com/mbvisti",
	"https://www.linkedin.com/in/mortenvistisen",
}

func SetSchema(builders ...SchemaBuilder) HeadDataOption {
	return func(hd *HeadData) {
		hd.SchemaBuilders = append(hd.SchemaBuilders, builders...)
	}
}

func buildStructuredData(data HeadData) map[string]any {
	if len(data.SchemaBuilders) == 0 {
		return nil
	}

	siteURL := buildCanonicalURL(data.canonical, "/")
	personID := siteURL + "#person"
	websiteID := siteURL + "#website"
	webpageID := data.canonicalURL + "#webpage"
	graph := []SchemaNode{
		{
			"@type":       "Person",
			"@id":         personID,
			"name":        "Morten Vistisen",
			"url":         siteURL,
			"image":       "https://media.mortenvistisen.com/mig_selv-removebg-preview.png",
			"email":       "blog@mortenvistisen.com",
			"jobTitle":    "Software Engineer",
			"description": "Danish software engineer and writer living in Spain, focused on Go, distributed systems, and practical software.",
			"sameAs":      personSameAs,
			"knowsAbout":  []string{"Go", "Distributed systems", "Web development", "Software architecture", "Bootstrapping"},
		},
		{
			"@type":       "WebSite",
			"@id":         websiteID,
			"name":        data.siteName,
			"url":         siteURL,
			"description": "Writing, projects, and field notes from software engineer Morten Vistisen.",
			"inLanguage":  "en-US",
			"author":      schemaRef(personID),
			"publisher":   schemaRef(personID),
		},
		{
			"@type":       "WebPage",
			"@id":         webpageID,
			"url":         data.canonicalURL,
			"name":        cleanPageTitle(data),
			"description": data.Description,
			"isPartOf":    schemaRef(websiteID),
			"author":      schemaRef(personID),
			"inLanguage":  "en-US",
		},
	}

	if data.Image != "" {
		graph[2]["primaryImageOfPage"] = SchemaNode{
			"@type": "ImageObject",
			"url":   buildCanonicalURL(data.canonical, data.Image),
		}
	}

	for _, builder := range data.SchemaBuilders {
		graph = mergeSchemaNodes(graph, builder(data)...)
	}

	return map[string]any{
		"@context": "https://schema.org",
		"@graph":   graph,
	}
}

func HomePageSchema() SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		return []SchemaNode{{
			"@type":      "ProfilePage",
			"@id":        data.canonicalURL + "#webpage",
			"mainEntity": schemaRef(personSchemaID(data)),
		}}
	}
}

func PageSchema(pageType string, breadcrumbs []SchemaBreadcrumb) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		return []SchemaNode{
			{
				"@type":      pageType,
				"@id":        data.canonicalURL + "#webpage",
				"mainEntity": schemaRef(personSchemaID(data)),
			},
			breadcrumbSchema(data, breadcrumbs),
		}
	}
}

func PostIndexSchema(items []models.ArticleEntity) SchemaBuilder {
	listItems := make([]schemaListItem, 0, len(items))
	for _, item := range items {
		listItem := schemaListItem{Name: item.Title, Path: routes.PostShow.URL(item.Slug)}
		if item.Excerpt.Valid {
			listItem.Description = item.Excerpt.String
		}
		if imageURL, ok := publicImageURL(item.ImageLink); ok {
			listItem.Image = imageURL
		}
		listItems = append(listItems, listItem)
	}
	return collectionSchema(
		"Writing",
		"Articles on Go, practical systems, products, and engineering.",
		[]SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Posts", Path: routes.PostIndex.URL()}},
		listItems,
	)
}

func PostShowSchema(post models.ArticleEntity, tags []models.TagEntity) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		description := ""
		if post.MetaDescription.Valid {
			description = post.MetaDescription.String
		} else if post.Excerpt.Valid {
			description = post.Excerpt.String
		}
		node := SchemaNode{
			"@type":            "BlogPosting",
			"@id":              data.canonicalURL + "#article",
			"headline":         post.Title,
			"description":      description,
			"url":              data.canonicalURL,
			"mainEntityOfPage": schemaRef(data.canonicalURL + "#webpage"),
			"author":           schemaRef(personSchemaID(data)),
			"publisher":        schemaRef(personSchemaID(data)),
			"keywords":         tagNames(tags),
		}
		if imageURL, ok := publicImageURL(post.ImageLink); ok {
			node["image"] = buildCanonicalURL(data.canonical, imageURL)
		}
		addPublishedDates(node, post.FirstPublishedAt.Valid, post.FirstPublishedAt.Time, post.UpdatedAt)
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Posts", Path: routes.PostIndex.URL()}, {Name: post.Title, Path: routes.PostShow.URL(post.Slug)}}),
			node,
		}
	}
}

func NewsletterIndexSchema(items []models.NewsletterEntity) SchemaBuilder {
	listItems := make([]schemaListItem, 0, len(items))
	for _, item := range items {
		listItem := schemaListItem{Name: item.Title, Path: routes.NewsletterShow.URL(item.Slug.String), Description: item.MetaDescription}
		if imageURL, ok := publicImageURL(item.ImageLink); ok {
			listItem.Image = imageURL
		}
		listItems = append(listItems, listItem)
	}
	return collectionSchema(
		"Newsletters",
		"Field notes and dispatches on software, products, and independent work.",
		[]SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Newsletters", Path: routes.NewsletterIndex.URL()}},
		listItems,
	)
}

func NewsletterShowSchema(newsletter models.NewsletterEntity) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		node := SchemaNode{
			"@type":            "BlogPosting",
			"@id":              data.canonicalURL + "#article",
			"headline":         newsletter.Title,
			"description":      newsletter.MetaDescription,
			"url":              data.canonicalURL,
			"mainEntityOfPage": schemaRef(data.canonicalURL + "#webpage"),
			"author":           schemaRef(personSchemaID(data)),
			"publisher":        schemaRef(personSchemaID(data)),
		}
		if imageURL, ok := publicImageURL(newsletter.ImageLink); ok {
			node["image"] = buildCanonicalURL(data.canonical, imageURL)
		}
		addPublishedDates(node, newsletter.ReleasedAt.Valid, newsletter.ReleasedAt.Time, newsletter.UpdatedAt)
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Newsletters", Path: routes.NewsletterIndex.URL()}, {Name: newsletter.Title, Path: routes.NewsletterShow.URL(newsletter.Slug.String)}}),
			node,
		}
	}
}

func ProjectIndexSchema(items []models.ProjectEntity) SchemaBuilder {
	listItems := make([]schemaListItem, 0, len(items))
	for _, item := range items {
		listItem := schemaListItem{Name: item.Title, Path: routes.ProjectShow.URL(item.Slug), Description: item.Description}
		if imageURL, ok := publicImageURL(item.ImageLink); ok {
			listItem.Image = imageURL
		}
		listItems = append(listItems, listItem)
	}
	return collectionSchema(
		"Projects",
		"Products and experiments built by Morten Vistisen.",
		[]SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Projects", Path: routes.ProjectIndex.URL()}},
		listItems,
	)
}

func ProjectShowSchema(project models.ProjectEntity) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		node := SchemaNode{
			"@type":       "CreativeWork",
			"@id":         data.canonicalURL + "#project",
			"name":        project.Title,
			"description": project.Description,
			"url":         data.canonicalURL,
			"creator":     schemaRef(personSchemaID(data)),
			"publisher":   schemaRef(personSchemaID(data)),
		}
		if imageURL, ok := publicImageURL(project.ImageLink); ok {
			node["image"] = buildCanonicalURL(data.canonical, imageURL)
		}
		if project.ProjectUrl.Valid {
			node["sameAs"] = project.ProjectUrl.String
		}
		if project.StartedAt.Valid {
			node["dateCreated"] = project.StartedAt.Time.Format(time.RFC3339)
		}
		if !project.UpdatedAt.IsZero() {
			node["dateModified"] = project.UpdatedAt.Format(time.RFC3339)
		}
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Projects", Path: routes.ProjectIndex.URL()}, {Name: project.Title, Path: routes.ProjectShow.URL(project.Slug)}}),
			node,
		}
	}
}

func collectionSchema(name, description string, breadcrumbs []SchemaBreadcrumb, items []schemaListItem) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		nodes := []SchemaNode{
			breadcrumbSchema(data, breadcrumbs),
			{
				"@type":       "CollectionPage",
				"@id":         data.canonicalURL + "#webpage",
				"name":        name,
				"description": description,
				"url":         data.canonicalURL,
			},
		}
		if len(items) > 0 {
			nodes = append(nodes, itemListSchema(data, name, items))
		}
		return nodes
	}
}

func breadcrumbSchema(data HeadData, breadcrumbs []SchemaBreadcrumb) SchemaNode {
	items := make([]SchemaNode, 0, len(breadcrumbs))
	for index, breadcrumb := range breadcrumbs {
		items = append(items, SchemaNode{
			"@type":    "ListItem",
			"position": index + 1,
			"name":     breadcrumb.Name,
			"item":     buildCanonicalURL(data.canonical, breadcrumb.Path),
		})
	}
	return SchemaNode{
		"@type":           "BreadcrumbList",
		"@id":             data.canonicalURL + "#breadcrumb",
		"itemListElement": items,
	}
}

func itemListSchema(data HeadData, name string, items []schemaListItem) SchemaNode {
	elements := make([]SchemaNode, 0, len(items))
	for index, item := range items {
		node := SchemaNode{
			"@type":    "ListItem",
			"position": index + 1,
			"url":      buildCanonicalURL(data.canonical, item.Path),
			"name":     item.Name,
		}
		if item.Description != "" {
			node["description"] = item.Description
		}
		if item.Image != "" {
			node["image"] = buildCanonicalURL(data.canonical, item.Image)
		}
		elements = append(elements, node)
	}
	return SchemaNode{
		"@type":           "ItemList",
		"@id":             data.canonicalURL + "#item-list",
		"name":            name,
		"itemListElement": elements,
	}
}

func tagNames(tags []models.TagEntity) []string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Title)
	}
	return names
}

func addPublishedDates(node SchemaNode, hasPublishedAt bool, publishedAt, updatedAt time.Time) {
	if hasPublishedAt {
		node["datePublished"] = publishedAt.Format(time.RFC3339)
	}
	if !updatedAt.IsZero() {
		node["dateModified"] = updatedAt.Format(time.RFC3339)
	}
}

func cleanPageTitle(data HeadData) string {
	return strings.TrimSuffix(data.Title, " - "+data.siteName)
}

func personSchemaID(data HeadData) string {
	return buildCanonicalURL(data.canonical, "/") + "#person"
}

func schemaRef(id string) SchemaNode {
	return SchemaNode{"@id": id}
}

func mergeSchemaNodes(graph []SchemaNode, nodes ...SchemaNode) []SchemaNode {
	for _, node := range nodes {
		id, ok := node["@id"].(string)
		if !ok || id == "" {
			graph = append(graph, node)
			continue
		}
		merged := false
		for i := range graph {
			if graphID, ok := graph[i]["@id"].(string); ok && graphID == id {
				maps.Copy(graph[i], node)
				merged = true
				break
			}
		}
		if !merged {
			graph = append(graph, node)
		}
	}
	return graph
}
