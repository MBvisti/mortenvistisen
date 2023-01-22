use actix_web::{Responder, web, HttpResponse, post};

#[post("/subscribe")]
pub async fn subscribe(
    tmpl: web::Data<tera::Tera>,
) -> impl Responder {
    let context = tera::Context::new();

    match tmpl.render("post.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            let error_page = tmpl.render("error_page.html", &context).unwrap();
            HttpResponse::InternalServerError()
                .content_type("text/html")
                .body(error_page)
        }
    }
}

#[post("/subscribe/verify")]
pub async fn verify_subscription(
    tmpl: web::Data<tera::Tera>,
) -> impl Responder {
    let context = tera::Context::new();

    match tmpl.render("post.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            let error_page = tmpl.render("error_page.html", &context).unwrap();
            HttpResponse::InternalServerError()
                .content_type("text/html")
                .body(error_page)
        }
    }
}
