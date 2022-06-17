use actix_web::{get, Responder, HttpResponse, web};

#[get("/")]
pub async fn index(tmpl: web::Data<tera::Tera>) -> impl Responder {
    println!("hitting index");
    let context = tera::Context::new();

    match tmpl.render("index.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{}", e);
            HttpResponse::Ok().body("Hey there!")
        },
    }
}

