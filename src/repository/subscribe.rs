use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::Email;

pub async fn create_new_subscriber(
    pool: &PgPool,
    new_subscriber: &Email,
    referer: &str,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO subscriptions (id, email, referer, subscribed_at, is_verified)
            VALUES ($1, $2, $3, $4, $5)
        "#,
        Uuid::new_v4(),
        new_subscriber.as_ref(),
        referer,
        Utc::now(),
        false,
    )
    .execute(pool)
    .await?;
    Ok(())
}

pub async fn does_email_exists(
    pool: &PgPool,
    email: &Email,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<bool, sqlx::Error> {
    match sqlx::query!(
        r#"
            select EXISTS(select email from subscriptions where email = $1) as "exists"
        "#,
        email.as_ref(),
    )
    .fetch_one(pool)
    .await {
        Ok(record) => {
            let exists = match record.exists {
                Some(e) => e,
                None => {
                    println!("exists returned null for email: {}", email.as_ref());
                    false
                },
            };
           return Ok(exists)
        },
        Err(e) => {
            return Err(e)
        },
    };
}

pub async fn is_user_verified(
    pool: &PgPool,
    email: &Email,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<bool, sqlx::Error> {
    match sqlx::query!(
        r#"
            select is_verified from subscriptions where email = $1
        "#,
        email.as_ref(),
    )
    .fetch_one(pool)
    .await {
        Ok(record) => {
           return Ok(record.is_verified)
        },
        Err(e) => {
            return Err(e)
        },
    };
}

pub async fn update_email_to_verified(
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

