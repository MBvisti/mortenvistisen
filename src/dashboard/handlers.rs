use actix_identity::Identity;
use actix_web::http::header::LOCATION;
use actix_web::{get, HttpResponse, Responder};

use crate::{
    auth_stuff::validate_cookie_identity,
    template::{render_internal_error_tmpl, render_template},
};

#[get("/dashboard")]
pub async fn index(identity: Option<Identity>) -> impl Responder {
    let _id = match validate_cookie_identity(identity) {
        Ok(id) => id,
        Err(_e) => {
            return HttpResponse::SeeOther()
                .insert_header((LOCATION, "/login"))
                .finish()
        }
    };
    let context = tera::Context::new();

    let tmpl = match render_template("dashboard_home.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}
