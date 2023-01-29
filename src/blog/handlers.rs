use actix_web::{get, HttpResponse, Responder};

use crate::{
    article::find_all_front_matter,
    template::{render_internal_error_tmpl, render_template},
};

#[tracing::instrument(name = "visit home page")]
#[get("/")]
pub async fn index() -> impl Responder {
    let mut context = tera::Context::new();

    let mut front_matters = match find_all_front_matter() {
        Ok(fm) => fm,
        Err(e) => {
            tracing::error!("failed to find all frontmatters: {:?}", e);

            return HttpResponse::InternalServerError()
                .content_type("text/html")
                .body(render_internal_error_tmpl(None));
        }
    };

    front_matters.sort_by(|a, b| b.order.cmp(&a.order));
    context.insert("posts", &front_matters);

    let tmpl = match render_template("home.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}
