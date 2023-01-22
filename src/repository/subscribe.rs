use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::Email;

pub async fn create_new_subscriber(
    pool: &PgPool,
    new_subscriber: &Email,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO subscriptions (id, email, subscribed_at, is_verified)
            VALUES ($1, $2, $3, $4)
        "#,
        Uuid::new_v4(),
        new_subscriber.as_ref(),
        Utc::now(),
        false,
    )
    .execute(pool)
    .await?;
    Ok(())
}

pub async fn update_email_to_verfied(
    pool: &PgPool,
    email: &Email,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            UPDATE subscriptions set is_verified = $2 where email = $1
        "#,
        email.as_ref(),
        true,
    )
    .execute(pool)
    .await?;
    Ok(())
}

