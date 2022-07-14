use actix_files as fs;
use handlers::not_found;
use std::net::TcpListener;

#[macro_use]
extern crate lazy_static;

use actix_web::{dev::Server, middleware, web, App, HttpResponse, HttpServer};
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

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
    let srv = HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(TEMPLATES.clone()))
            .wrap(middleware::Logger::default()) // enable logger
            .route("/status", web::get().to(HttpResponse::Ok))
            .service(fs::Files::new("/static", "static/").use_last_modified(true))
            .service(handlers::index)
            .service(handlers::render_post)
            .default_service(web::route().to(not_found))
    })
    .listen(listener)?
    .run();
    Ok(srv)
}
