-- +goose Up
-- +goose StatementBegin
UPDATE tags
SET title = BTRIM(title);

UPDATE tags
SET title = 'tag-' || id
WHERE title = '';

WITH canonical_tags AS (
    SELECT LOWER(title) AS normalized_title, MIN(id) AS canonical_id
    FROM tags
    GROUP BY LOWER(title)
    HAVING COUNT(*) > 1
)
UPDATE article_tag_connections AS connection
SET tag_id = canonical_tags.canonical_id
FROM tags
JOIN canonical_tags ON canonical_tags.normalized_title = LOWER(tags.title)
WHERE connection.tag_id = tags.id
  AND tags.id <> canonical_tags.canonical_id;

DELETE FROM article_tag_connections AS connection
USING article_tag_connections AS duplicate
WHERE connection.article_id = duplicate.article_id
  AND connection.tag_id = duplicate.tag_id
  AND connection.id > duplicate.id;

DELETE FROM tags AS tag
USING tags AS duplicate
WHERE LOWER(tag.title) = LOWER(duplicate.title)
  AND tag.id > duplicate.id;

ALTER TABLE article_tag_connections
    DROP CONSTRAINT article_tag_connections_article_id_fkey,
    DROP CONSTRAINT article_tag_connections_tag_id_fkey,
    ADD CONSTRAINT article_tag_connections_article_id_fkey
        FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    ADD CONSTRAINT article_tag_connections_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE;

ALTER TABLE tags
    ADD CONSTRAINT tags_title_not_blank CHECK (BTRIM(title) <> '');

CREATE UNIQUE INDEX tags_normalized_title_unique_idx
    ON tags (LOWER(BTRIM(title)));

CREATE UNIQUE INDEX article_tag_connections_article_tag_unique_idx
    ON article_tag_connections (article_id, tag_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS article_tag_connections_article_tag_unique_idx;
DROP INDEX IF EXISTS tags_normalized_title_unique_idx;

ALTER TABLE tags
    DROP CONSTRAINT IF EXISTS tags_title_not_blank;

ALTER TABLE article_tag_connections
    DROP CONSTRAINT article_tag_connections_article_id_fkey,
    DROP CONSTRAINT article_tag_connections_tag_id_fkey,
    ADD CONSTRAINT article_tag_connections_article_id_fkey
        FOREIGN KEY (article_id) REFERENCES articles(id),
    ADD CONSTRAINT article_tag_connections_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES tags(id);
-- +goose StatementEnd
