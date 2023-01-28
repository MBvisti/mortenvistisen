use secrecy::ExposeSecret;
use std::net::TcpListener;
use tracing::dispatcher::set_global_default;
use tracing_log::LogTracer;

use actix_web::web;
use tracing_bunyan_formatter::{BunyanFormattingLayer, JsonStorageLayer};
use tracing_subscriber::{layer::SubscriberExt, EnvFilter, Registry};

use mortenvistisen_blog::{
    configuration::get_config, email_client::EmailClient, start_blog, subscriber::Email,
};
use sqlx::postgres::PgPoolOptions;

#[tokio::main]
async fn main() -> std::io::Result<()> {
    println!("Starting yo");
    let cfg = match get_config() {
        Ok(cfg) => cfg,
        Err(e) => panic!("{}", e),
    };
    let db_conn_pool = match PgPoolOptions::new()
        .idle_timeout(std::time::Duration::from_secs(2))
        .acquire_timeout(std::time::Duration::from_secs(2))
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

    let _guard = sentry::init((
        cfg.server.sentry_dns.expose_secret().to_string(),
        sentry::ClientOptions {
            traces_sample_rate: 0.2,
            release: sentry::release_name!(),
            max_breadcrumbs: 50,
            debug: false, // <- this should only be used during development
            ..Default::default()
        },
    ));

    LogTracer::init().expect("failed to set logger");

    let env_filter = EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info"));
    let log_formatting_layer =
        BunyanFormattingLayer::new("mortenvistisen_blog".into(), std::io::stdout);

    let log_subscriber = Registry::default()
        .with(env_filter)
        .with(JsonStorageLayer)
        .with(log_formatting_layer)
        .with(sentry_tracing::layer());

    set_global_default(log_subscriber.into()).expect("failed to set log subscriber");

    let listener = match TcpListener::bind(format!("{}:{}", cfg.server.host, cfg.server.port)) {
        Ok(l) => l,
        Err(e) => panic!("{}", e),
    };

    start_blog(listener, db_conn_pool, web::Data::new(email_client))?.await?;
    Ok(())
}
