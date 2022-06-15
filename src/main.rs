use actix_web::{http::Error, web, App, HttpResponse, HttpServer, error, middleware};
use tera::{Context, Tera};

// store tera template in application state
async fn index(tmpl: web::Data<tera::Tera>) -> Result<HttpResponse, Error> {
    println!("hitting index");
    let context = Context::new();

    let s = match tmpl
        .render("index.html", &context) {
            Ok(s) => s,
            Err(e) => {
                println!("{}", e);
                error::ErrorInternalServerError(e).to_string()
            },
        };
    Ok(HttpResponse::Ok().content_type("text/html").body(s))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    std::env::set_var("RUST_LOG", "actix_web=debug");
    env_logger::init();

    println!("Listening on: 127.0.0.1:8080, open browser and visit have a try!");
    HttpServer::new(|| {
        let tera = match Tera::new("templates/**/*.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {}", e);
                ::std::process::exit(1);
            }
        };
        App::new()
            .app_data(web::Data::new(tera))
            .wrap(middleware::Logger::default()) // enable logger
            .service(web::resource("/").route(web::get().to(index)))
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
