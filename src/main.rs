use actix_web::web;
use mortenvistisen_blog::{
    auth_stuff::create_hashed_password_from_string,
    configuration::get_config,
    email_client::EmailClient,
    repository::{create_new_user, does_user_exists},
    start_blog,
    subscriber::Email,
    telemetry::{get_subscriber, init_subscriber},
};
use secrecy::{ExposeSecret, Secret};
use sqlx::{postgres::PgPoolOptions, PgPool};
use std::net::TcpListener;

async fn ensure_admin_user_exists(
    pool: &PgPool,
    salty: Secret<String>,
    password: Secret<String>,
) {
    let email = Email::parse("mbv@mortenvistisen.com".to_string()).unwrap();
    let exsits = match does_user_exists(&pool, &email).await {
        Ok(exists) => exists,
        Err(e) => panic!("{}", e),
    };

    if !exsits {
        let hp =
            create_hashed_password_from_string(password.expose_secret(), salty.expose_secret());
        if let Err(create_err) = create_new_user(&pool, &email, &hp).await {
            panic!("{}", create_err)
        }
    }
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let cfg = match get_config() {
        Ok(cfg) => cfg,
        Err(e) => panic!("{}", e),
    };

    let subscriber = get_subscriber(
        "mortenvistisen_blog".into(),
        "error".into(),
        std::io::stdout,
        cfg.server.sentry_dns,
    );
    init_subscriber(subscriber);

    let db_conn_pool = match PgPoolOptions::new()
        .idle_timeout(std::time::Duration::from_secs(10))
        .acquire_timeout(std::time::Duration::from_secs(10))
        .connect_lazy(&cfg.database.get_connection_string())
    {
        Ok(db_conn) => db_conn,
        Err(e) => panic!("{}", e),
    };

    ensure_admin_user_exists(&db_conn_pool, cfg.server.password_salt, Secret::new("yEkh9zS85DD7xcMmB9li".to_string())).await;

    let sender = match Email::parse(cfg.email_client.sender) {
        Ok(s) => s,
        Err(e) => panic!("{}", e),
    };
    let email_client = EmailClient::new(
        cfg.email_client.api_base_url,
        cfg.server.base_url.clone(),
        sender,
        cfg.email_client.auth_token,
    );

    let listener = match TcpListener::bind(format!("{}:{}", cfg.server.host, cfg.server.port)) {
        Ok(l) => l,
        Err(e) => panic!("{}", e),
    };

    start_blog(listener, db_conn_pool, web::Data::new(email_client))?.await?;
    Ok(())
}
