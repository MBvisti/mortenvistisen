use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::Email;

pub async fn create_new_subscriber(
    pool: &PgPool,
    new_subscriber: &Email,
    referer: &str,
    subscriber_id: Uuid,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO subscriptions (id, email, referer, subscribed_at, is_verified)
            VALUES ($1, $2, $3, $4, $5)
        "#,
        subscriber_id,
        new_subscriber.as_ref(),
        referer,
        Utc::now(),
        false,
    )
    .execute(pool)
    .await?;
    Ok(())
}

pub async fn create_new_subscriber_token(
    pool: &PgPool,
    subscriber_id: Uuid,
    token: &str,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO subscription_token (subscription_token, subscription_id)
            VALUES ($1, $2)     
        "#,
        token,
        subscriber_id,
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

pub async fn is_user_verified(
    pool: &PgPool,
    subscriber_id: Uuid,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<bool, sqlx::Error> {
    match sqlx::query!(
        r#"
            select is_verified from subscriptions where id = $1
        "#,
        subscriber_id,
    )
    .fetch_one(pool)
    .await
    {
        Ok(record) => Ok(record.is_verified),
        Err(e) => Err(e),
    }
}

pub async fn get_subscriber_id_from_token(
    pool: &PgPool,
    tkn: &str,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<Uuid, sqlx::Error> {
    match sqlx::query!(
        r#"
            select subscription_id from subscription_token where subscription_token = $1
        "#,
        tkn,
    )
    .fetch_one(pool)
    .await
    {
        Ok(record) => Ok(record.subscription_id),
        Err(e) => Err(e),
    }
}

pub async fn update_email_to_verified(
    pool: &PgPool,
    id: Uuid,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            UPDATE subscriptions set is_verified = $2 where id = $1
        "#,
        id,
        true,
    )
    .execute(pool)
    .await?;
    Ok(())
}

pub async fn delete_subscriber(
    pool: &PgPool,
    id: Uuid,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            DELETE from subscriptions  where id = $1
        "#,
        id,
    )
    .execute(pool)
    .await?;
    Ok(())
}

pub async fn delete_subscriber_token(
    pool: &PgPool,
    tkn: &str,
    // using the () type, known as unit type, and has the purpose of being useless. I.e.
    // we don't have anything meaningful to return
) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            DELETE from subscription_token  where subscription_token = $1
        "#,
        tkn,
    )
    .execute(pool)
    .await?;
    Ok(())
}
