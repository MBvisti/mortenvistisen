use std::{net::TcpListener};

#[macro_use]
extern crate lazy_static;

use actix_web::{dev::Server, HttpServer, middleware, App, web};
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
        // tera.register_filter("do_nothing", do_nothing_filter);
        tera
    };
}

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
    let srv = HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(TEMPLATES.clone()))
            .wrap(middleware::Logger::default()) // enable logger
            .service(handlers::index)
    })
    .listen(listener)?
    .run();
    Ok(srv)
}
