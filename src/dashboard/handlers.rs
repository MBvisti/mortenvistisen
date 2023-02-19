use actix_web::{get, Responder, HttpResponse};

use crate::template::{render_internal_error_tmpl, render_template};


#[get("/dashboard")]
pub async fn index() -> impl Responder {
    let context = tera::Context::new();

    let tmpl = match render_template("dashboard_home.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}
