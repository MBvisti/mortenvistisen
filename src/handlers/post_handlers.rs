use std::fs;

use actix_web::{get, web, HttpResponse, Responder};
use pulldown_cmark::{html, Options, Parser};

use crate::{
    handlers::FrontMatter,
    template::{render_not_found_error_tmpl, render_template, render_internal_error_tmpl},
};

#[get("/posts/{post_name}")]
pub async fn render_post(post_name: web::Path<String>) -> impl Responder {
    let mut context = tera::Context::new();

    let options = Options::empty();
    let markdown_input = match fs::read_to_string(format!("./posts/{}/article.md", post_name)) {
        Ok(s) => s,
        Err(e) => {
            println!("{:?}", e);

            return HttpResponse::NotFound()
                .content_type("text/html")
                .body(render_not_found_error_tmpl(None));
        }
    };

    let front_matter_input =
        match fs::read_to_string(format!("./posts/{}/article_frontmatter.toml", post_name)) {
            Ok(s) => s,
            Err(e) => {
                println!("{:?}", e);

                return HttpResponse::NotFound()
                    .content_type("text/html")
                    .body(render_not_found_error_tmpl(None));
            }
        };

    let front_matter: FrontMatter = match toml::from_str(&front_matter_input) {
        Ok(fm) => fm,
        Err(e) => {
            println!("{:?}", e);

            return HttpResponse::NotFound()
                .content_type("text/html")
                .body(render_not_found_error_tmpl(None));
        }
    };

    let parser = Parser::new_ext(&markdown_input, options);

    let mut html_output = String::new();
    html::push_html(&mut html_output, parser);

    context.insert("post", &html_output);
    context.insert("meta_data", &front_matter);

    let tmpl = match render_template("post.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None)
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}
