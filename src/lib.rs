use actix_files as fs;
use email_client::EmailClient;
use sqlx::PgPool;
use std::net::TcpListener;

#[macro_use]
extern crate lazy_static;

use actix_web::{
    dev::Server, get, middleware, web, App, Error, HttpRequest, HttpResponse, HttpServer, Responder,
};

pub mod article;
pub mod blog;
pub mod configuration;
pub mod email_client;
pub mod repository;
pub mod subscriber;
pub mod template;

#[get("/robots.txt")]
async fn robots_text(_req: HttpRequest) -> Result<fs::NamedFile, Error> {
    let file = fs::NamedFile::open_async("static/robots.txt").await?;
    Ok(file.use_last_modified(true))
}

#[get("/sitemap.xml")]
async fn sitemap_text(_req: HttpRequest) -> Result<fs::NamedFile, Error> {
    let file = fs::NamedFile::open_async("static/sitemap.xml").await?;
    Ok(file.use_last_modified(true))
}

async fn not_found(tmpl: web::Data<tera::Tera>) -> impl Responder {
    let not_found_page = tmpl
        .render("not_found.html", &tera::Context::new())
        .unwrap();
    HttpResponse::InternalServerError()
        .content_type("text/html")
        .body(not_found_page)
}

pub fn start_blog(
    listener: TcpListener,
    db_pool: PgPool,
    email_client: web::Data<EmailClient>,
) -> Result<Server, std::io::Error> {
    let db_conn_pool = web::Data::new(db_pool);
    let srv = HttpServer::new(move || {
        App::new()
            .app_data(db_conn_pool.clone())
            .app_data(email_client.clone())
            .wrap(middleware::Logger::default()) // enable logger
            .route("/status", web::get().to(HttpResponse::Ok))
            .service(robots_text)
            .service(sitemap_text)
            .service(fs::Files::new("/static", "static/").use_last_modified(true))
            // .service(fs::Files::new("/static", "static/robots.txt").use_last_modified(true))
            .service(blog::index)
            .service(article::render_post)
            .service(subscriber::subscribe)
            .service(subscriber::verify_subscription)
            .service(subscriber::delete_subscriber)
            .default_service(web::route().to(not_found))
    })
    .listen(listener)?
    .run();
    Ok(srv)
}
