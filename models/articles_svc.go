package models

// type postDatabase interface {
// 	InsertPost(ctx context.Context, arg database.InsertPostParams) (uuid.UUID, error)
// 	UpdatePost(ctx context.Context, arg database.UpdatePostParams) (database.Post, error)
// 	AssociateTagWithPost(ctx context.Context, arg database.AssociateTagWithPostParams) error
// }
//
// func NewArticle(
// 	ctx context.Context,
// 	db postDatabase,
// 	v *validator.Validate,
// 	newArticle domain.NewArticle,
// 	associatedTags []string,
// ) error {
// 	if err := v.Struct(newArticle); err != nil {
// 		telemetry.Logger.Error("provided post data did not pass the validation", "error", err)
// 		return err
// 	}
//
// 	now := time.Now()
//
// 	args := database.InsertPostParams{
// 		ID:          uuid.New(),
// 		CreatedAt:   database.ConvertToPGTimestamp(now),
// 		UpdatedAt:   database.ConvertToPGTimestamp(now),
// 		Title:       newArticle.Title,
// 		HeaderTitle: sql.NullString{Valid: true, String: newArticle.HeaderTitle},
// 		Filename:    newArticle.Filename,
// 		Slug:        slug.MakeLang(newArticle.Title, "en"),
// 		Excerpt:     newArticle.Excerpt,
// 		Draft:       true,
// 	}
//
// 	if newArticle.ReleaseNow {
// 		args.ReleasedAt = database.ConvertToPGTimestamp(now)
// 		args.Draft = false
// 	}
//
// 	id, err := db.InsertPost(ctx, args)
// 	if err != nil {
// 		return err
// 	}
//
// 	// TODO: run in transaction
// 	for _, associatedTag := range associatedTags {
// 		tagID, err := uuid.Parse(associatedTag)
// 		if err != nil {
// 			return err
// 		}
//
// 		if err := db.AssociateTagWithPost(
// 			ctx,
// 			database.AssociateTagWithPostParams{
// 				ID:     uuid.New(),
// 				PostID: id,
// 				TagID:  tagID,
// 			}); err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// func UpdateArticle(
// 	ctx context.Context,
// 	db postDatabase,
// 	v *validator.Validate,
// 	updateArticle domain.UpdateArticle,
// ) (domain.Article, error) {
// 	if err := v.Struct(updateArticle); err != nil {
// 		telemetry.Logger.Error("provided post data did not pass the validation", "error", err)
// 		return domain.Article{}, err
// 	}
//
// 	now := time.Now()
// 	args := database.UpdatePostParams{
// 		ID:          updateArticle.ID,
// 		UpdatedAt:   database.ConvertToPGTimestamp(now),
// 		Title:       updateArticle.Title,
// 		HeaderTitle: sql.NullString{Valid: true, String: updateArticle.HeaderTitle},
// 		Slug:        slug.MakeLang(updateArticle.Title, "en"),
// 		Excerpt:     updateArticle.Excerpt,
// 		Draft:       updateArticle.ReleaseNow,
// 		ReleasedAt:  database.ConvertToPGTimestamp(now),
// 		ReadTime:    sql.NullInt32{Int32: updateArticle.EstimatedReadTime, Valid: true},
// 	}
//
// 	updatedArticle, err := db.UpdatePost(ctx, args)
// 	if err != nil {
// 		return domain.Article{}, err
// 	}
//
// 	return domain.Article{
// 		ID:          updatedArticle.ID,
// 		CreatedAt:   database.ConvertFromPGTimestampToTime(updatedArticle.CreatedAt),
// 		UpdatedAt:   database.ConvertFromPGTimestampToTime(updatedArticle.UpdatedAt),
// 		Title:       updatedArticle.Title,
// 		Filename:    updatedArticle.Filename,
// 		Slug:        updatedArticle.Slug,
// 		Excerpt:     updatedArticle.Excerpt,
// 		Draft:       updatedArticle.Draft,
// 		ReleaseDate: database.ConvertFromPGTimestampToTime(updatedArticle.ReleasedAt),
// 		ReadTime:    updatedArticle.ReadTime.Int32,
// 	}, nil
// }
