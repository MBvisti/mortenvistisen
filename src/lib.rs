use actix_files as fs;
use handlers::not_found;
use std::net::TcpListener;

#[macro_use]
extern crate lazy_static;

use actix_web::{
    dev::Server, get, middleware, web, App, Error, HttpRequest, HttpResponse, HttpServer,
};
use tera::Tera;

pub mod handlers;

lazy_static! {
    pub static ref TEMPLATES: Tera = {
        let mut tera = match Tera::new("templates/**/*.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {}", e);
                ::std::process::exit(1);
            }
        };
        tera.autoescape_on(vec![".html", ".sql"]);
        tera
    };
}

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

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
    let tmpl = web::Data::new(TEMPLATES.clone());
    let srv = HttpServer::new(move || {
        App::new()
            .app_data(web::Data::clone(&tmpl))
            .wrap(middleware::Logger::default()) // enable logger
            .route("/status", web::get().to(HttpResponse::Ok))
            .service(robots_text)
            .service(sitemap_text)
            .service(fs::Files::new("/static", "static/").use_last_modified(true))
            // .service(fs::Files::new("/static", "static/robots.txt").use_last_modified(true))
            .service(handlers::index)
            .service(handlers::render_post)
            .default_service(web::route().to(not_found))
    })
    .listen(listener)?
    .run();
    Ok(srv)
}
