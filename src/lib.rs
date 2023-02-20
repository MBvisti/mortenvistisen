use actix_files as fs;
use actix_identity::{Identity, IdentityMiddleware};
use actix_session::{config::PersistentSession, storage::CookieSessionStore, SessionMiddleware};
use email_client::EmailClient;
use reqwest::StatusCode;
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use std::net::TcpListener;
use tracing_actix_web::TracingLogger;

#[macro_use]
extern crate lazy_static;

use actix_web::{
    cookie::{time::Duration, Key},
    dev::Server,
    get, post, web, App, Error, HttpMessage, HttpRequest, HttpResponse, HttpServer, Responder,
};

use crate::{
    auth_stuff::verify_hashed_password,
    repository::get_user_hashed_password,
    subscriber::Email,
    template::{render_internal_error_tmpl, render_template},
};

pub mod article;
pub mod auth_stuff;
pub mod blog;
pub mod configuration;
pub mod dashboard;
pub mod email_client;
pub mod repository;
pub mod subscriber;
pub mod telemetry;
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

#[derive(Debug, Serialize)]
struct LoginMetaData {
    error_msg: Option<String>,
    has_error: bool,
    is_success: bool,
}
#[get("/login")]
async fn login_handler(req: HttpRequest) -> impl Responder {
    let mut context = tera::Context::new();
    context.insert(
        "meta_data",
        &LoginMetaData {
            error_msg: None,
            has_error: false,
            is_success: true,
        },
    );

    Identity::login(&req.extensions(), "user1".to_owned()).unwrap();
    let tmpl = match render_template("login.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}

#[get("/auth-redirect")]
async fn auth_redirect_handler(_req: HttpRequest) -> impl Responder {
    HttpResponse::Ok()
        .append_header(("HX-Redirect", "/dashboard"))
        .finish()
}

#[derive(Debug, Deserialize)]
pub struct LoginFormData {
    pub email: String,
    pub password: String,
}
#[post("/login")]
async fn authenticate_handler(
    _req: HttpRequest,
    form: web::Form<LoginFormData>,

    pool: web::Data<PgPool>,
) -> impl Responder {
    let email = match Email::parse(form.email.to_string()) {
        Ok(email) => email,
        Err(e) => {
            tracing::error!("could not parse email bc: {:?}", e);
            return web::Redirect::to("/login").using_status_code(StatusCode::BAD_REQUEST);
        }
    };

    let user_hp = match get_user_hashed_password(&pool, &email).await {
        Ok(user_hp) => user_hp,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                return web::Redirect::to("/login").using_status_code(StatusCode::BAD_REQUEST)
            }
            _ => {
                tracing::error!("could not get user hashed password bc: {:?}", e);
                return web::Redirect::to("/login")
                    .using_status_code(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };

    match verify_hashed_password(&user_hp, &form.password) {
        Ok(is_verified) => {
            if is_verified {
                return web::Redirect::to("/dashboard").using_status_code(StatusCode::FOUND);
            } else {
                return web::Redirect::to("/login").using_status_code(StatusCode::UNAUTHORIZED);
            }
        }
        Err(e) => {
            tracing::error!("could not get user hashed password bc: {:?}", e);
            return web::Redirect::to("/login")
                .using_status_code(StatusCode::INTERNAL_SERVER_ERROR);
        }
    }
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
    let session_key = Key::generate();
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
            .service(robots_text)
            .service(sitemap_text)
            .service(fs::Files::new("/static", "static/").use_last_modified(true))
            // .service(fs::Files::new("/static", "static/robots.txt").use_last_modified(true))
            .service(login_handler)
            .service(authenticate_handler)
            .service(auth_redirect_handler)
            .service(blog::index)
            .service(article::render_post)
            .service(subscriber::subscribe)
            .service(subscriber::verify_subscription)
            .service(subscriber::delete_subscriber)
            .service(dashboard::index)
            .default_service(web::route().to(not_found))
    })
    .listen(listener)?
    .run();
    Ok(srv)
}
