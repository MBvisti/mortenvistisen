use actix_files as fs;
use actix_identity::IdentityMiddleware;
use actix_session::{config::PersistentSession, storage::CookieSessionStore, SessionMiddleware};
use email_client::EmailClient;
use sqlx::PgPool;
use std::net::TcpListener;
use tracing_actix_web::TracingLogger;

use actix_web::{
    cookie::{time::Duration, Key},
    dev::Server,
    get,
    web::{self, Data},
    App, Error, HttpRequest, HttpResponse, HttpServer,
};

#[macro_use]
extern crate lazy_static;

pub mod configuration;
pub mod controllers;
pub mod email_client;
pub mod entities;
pub mod repository;
pub mod services;
pub mod telemetry;
pub mod views;

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

#[get("/sitemap_index.xml")]
async fn sitemap_index_text(_req: HttpRequest) -> Result<fs::NamedFile, Error> {
    let file = fs::NamedFile::open_async("static/sitemapindex.xml").await?;
    Ok(file.use_last_modified(true))
}

// #[derive(Debug, Serialize)]
// struct LoginMetaData {
//     error_msg: Option<String>,
//     has_error: bool,
//     is_success: bool,
// }
// #[get("/login")]
// async fn login_handler(_req: HttpRequest) -> impl Responder {
//     let mut context = tera::Context::new();
//     context.insert(
//         "meta_data",
//         &LoginMetaData {
//             error_msg: None,
//             has_error: false,
//             is_success: false,
//         },
//     );

//     let tmpl = match render_template("login.html", &context) {
//         Ok(t) => t,
//         Err(_) => render_internal_error_tmpl(None),
//     };

//     HttpResponse::Ok().content_type("text/html").body(tmpl)
// }

// #[get("/auth-redirect")]
// async fn auth_redirect_handler(_req: HttpRequest) -> impl Responder {
//     HttpResponse::Ok()
//         .append_header(("HX-Redirect", "/dashboard"))
//         .finish()
// }

// #[derive(Debug, Deserialize)]
// pub struct LoginFormData {
//     pub email: String,
//     pub password: String,
// }
// #[post("/login")]
// async fn authenticate_handler(
//     req: HttpRequest,
//     form: web::Form<LoginFormData>,
//     pool: Data<PgPool>,
// ) -> impl Responder {
//     let mut context = tera::Context::new();

//     let email = match Email::parse(form.email.to_string()) {
//         Ok(email) => email,
//         Err(e) => {
//             tracing::error!("could not parse email bc: {:?}", e);

//             context.insert(
//                 "meta_data",
//                 &LoginMetaData {
//                     error_msg: None,
//                     has_error: true,
//                     is_success: false,
//                 },
//             );
//             let tmpl = match render_template("login.html", &context) {
//                 Ok(t) => t,
//                 Err(_) => render_internal_error_tmpl(None),
//             };

//             return HttpResponse::Ok().content_type("text/html").body(tmpl);
//         }
//     };

//     let user_hp = match get_user_hashed_password(&pool, &email).await {
//         Ok(user_hp) => user_hp,
//         Err(e) => match e {
//             sqlx::Error::RowNotFound => {
//                 context.insert(
//                     "meta_data",
//                     &LoginMetaData {
//                         error_msg: None,
//                         has_error: true,
//                         is_success: false,
//                     },
//                 );
//                 let tmpl = match render_template("login.html", &context) {
//                     Ok(t) => t,
//                     Err(_) => render_internal_error_tmpl(None),
//                 };

//                 return HttpResponse::Ok().content_type("text/html").body(tmpl);
//             }
//             _ => {
//                 tracing::error!("could not get user hashed password bc: {:?}", e);
//                 context.insert(
//                     "meta_data",
//                     &LoginMetaData {
//                         error_msg: None,
//                         has_error: true,
//                         is_success: false,
//                     },
//                 );
//                 let tmpl = match render_template("login.html", &context) {
//                     Ok(t) => t,
//                     Err(_) => render_internal_error_tmpl(None),
//                 };

//                 return HttpResponse::Ok().content_type("text/html").body(tmpl);
//             }
//         },
//     };

//     match verify_password(&form.password, &user_hp) {
//         true => {
//             context.insert(
//                 "meta_data",
//                 &LoginMetaData {
//                     error_msg: None,
//                     has_error: false,
//                     is_success: true,
//                 },
//             );
//             let tmpl = match render_template("login.html", &context) {
//                 Ok(t) => t,
//                 Err(_) => render_internal_error_tmpl(None),
//             };

//             Identity::login(&req.extensions(), "user1".to_owned()).unwrap();

//             return HttpResponse::Ok().content_type("text/html").body(tmpl);
//         }
//         false => {
//             context.insert(
//                 "meta_data",
//                 &LoginMetaData {
//                     error_msg: None,
//                     has_error: true,
//                     is_success: false,
//                 },
//             );
//             let tmpl = match render_template("login.html", &context) {
//                 Ok(t) => t,
//                 Err(_) => render_internal_error_tmpl(None),
//             };

//             return HttpResponse::Ok().content_type("text/html").body(tmpl);
//         }
//     }
// }

// async fn not_found(tmpl: Data<tera::Tera>) -> impl Responder {
//     let not_found_page = tmpl
//         .render("not_found.html", &tera::Context::new())
//         .unwrap();
//     HttpResponse::InternalServerError()
//         .content_type("text/html")
//         .body(not_found_page)
// }

pub fn start_blog(
    listener: TcpListener,
    db_pool: PgPool,
    email_client: Data<EmailClient>,
    session_key: Key,
) -> Result<Server, std::io::Error> {
    let db_conn_pool = Data::new(db_pool);
    let srv = HttpServer::new(move || {
        let session =
            SessionMiddleware::builder(CookieSessionStore::default(), session_key.clone())
                .cookie_name("mbv_auth".to_string())
                .cookie_secure(false)
                .session_lifecycle(
                    PersistentSession::default().session_ttl(Duration::seconds(60 * 60 * 24 * 7)),
                )
                .build();
        App::new()
            .app_data(db_conn_pool.clone())
            .app_data(email_client.clone())
            .wrap(TracingLogger::default()) // enable logger
            .wrap(IdentityMiddleware::default())
            .wrap(session)
            .route("/status", web::get().to(HttpResponse::Ok))
            .service(fs::Files::new("/static", "static/").use_last_modified(true))
            .service(robots_text)
            .service(sitemap_text)
            .service(sitemap_index_text)
            .service(controllers::home_index)
            .service(controllers::render_post)
            .service(controllers::subscribe_to_newsletter)
            .service(controllers::verify_subscription)
            .service(controllers::delete_subscriber)
        // .default_service(web::route().to(not_found))
    })
    .listen(listener)?
    .run();
    Ok(srv)
}
