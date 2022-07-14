use std::fs;

use actix_web::{get, web, HttpResponse, Responder};
use pulldown_cmark::{html, Options, Parser};

#[get("/posts/{post_name}")]
pub async fn render_post(
    tmpl: web::Data<tera::Tera>,
    post_name: web::Path<String>,
) -> impl Responder {
    let mut context = tera::Context::new();

    let options = Options::empty();

    let markdown_input =
        match fs::read_to_string(format!("./posts/{}/article.md", post_name.into_inner())) {
            Ok(s) => s,
            Err(e) => {
                println!("{:?}", e);
                let not_found_page = tmpl.render("404_page.html", &context).unwrap();
                return HttpResponse::NotFound()
                    .content_type("text/html")
                    .body(not_found_page);
            }
        };

    let parser = Parser::new_ext(&markdown_input, options);

    let mut html_output = String::new();
    html::push_html(&mut html_output, parser);

    context.insert("post", &html_output);
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
