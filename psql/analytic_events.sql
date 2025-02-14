-- name: InsertAnalyticEvent :exec
insert into analytic_events (
    id,
	analytics_id,
	element_tag,
	element_text,
	element_class,
	element_id,
	element_href,
	custom_data
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
);
