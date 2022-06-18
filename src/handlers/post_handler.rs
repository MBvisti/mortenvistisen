use std::fs;

use actix_web::{get, web, HttpResponse, Responder};
use pulldown_cmark::{Parser, html, Options};

#[get("/posts")]
pub async fn render_post(tmpl: web::Data<tera::Tera>) -> impl Responder {
    let mut context = tera::Context::new();

    let options = Options::empty();

    let markdown_input = match fs::read_to_string("./posts/practical-approach-to-structuring-go-apps/article.md") {
        Ok(s) => s,
        Err(_) => "hello".to_string(),
    };
    
    let parser = Parser::new_ext(&markdown_input, options);

    let mut html_output = String::new();
    html::push_html(&mut html_output, parser);

    context.insert("post", &html_output);
    match tmpl.render("post.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            HttpResponse::InternalServerError().content_type("text/html").body("")
        }
    }
}
