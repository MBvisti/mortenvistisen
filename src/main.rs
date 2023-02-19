use actix_web::web;
use mortenvistisen_blog::{
    configuration::get_config,
    email_client::EmailClient,
    start_blog,
    subscriber::Email,
    telemetry::{get_subscriber, init_subscriber},
};
use sqlx::postgres::PgPoolOptions;
use std::net::TcpListener;

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
