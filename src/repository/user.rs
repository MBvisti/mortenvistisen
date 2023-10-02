use crate::entities::Email;
use sqlx::PgPool;
use uuid::Uuid;

pub async fn create_new_user(
    pool: &PgPool,
    email: &Email,
    hashed_password: &str,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO "user" (id, email, hashed_password)
            VALUES ($1, $2, $3)
        "#,
        Uuid::new_v4(),
        email.as_ref(),
        hashed_password,
    )
    .execute(pool)
    .await?;
    Ok(())
}

pub async fn get_user_hashed_password(
    pool: &PgPool,
    email: &Email,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<String, sqlx::Error> {
    let record = sqlx::query!(
        r#"
            select hashed_password from "user" where email = $1;
        "#,
        email.as_ref(),
    )
    .fetch_one(pool)
    .await
    .map_err(|e| {
        tracing::error!("failed to get user hashed password, with err: {:?}", e);
        e
    })?;

    Ok(record.hashed_password)
}

pub async fn does_user_exists(
    pool: &PgPool,
    email: &Email,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<bool, sqlx::Error> {
    match sqlx::query!(
        r#"
            select EXISTS(select email from "user" where email = $1) as "exists"
        "#,
        email.as_ref(),
    )
    .fetch_one(pool)
    .await
    {
        Ok(record) => {
            let exists = match record.exists {
                Some(e) => e,
                None => {
                    println!("exists returned null for email: {}", email.as_ref());
                    false
                }
            };
            Ok(exists)
        }
        Err(e) => Err(e),
    }
}
