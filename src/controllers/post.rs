use crate::entities::FrontMatter;
use std::fs;

use actix_web::{get, web, Responder};
use pulldown_cmark::{html, Options, Parser};

use crate::{
    views::{render_internal_error_tmpl, render_not_found_error_tmpl, render_template},
    views::Post,
};

#[tracing::instrument(name = "visit article")]
#[get("/posts/{post_name}")]
pub async fn render_post(post_name: web::Path<String>) -> impl Responder {
    let options = Options::empty();
    let markdown_input = match fs::read_to_string(format!("./posts/{post_name}/article.md")) {
        Ok(s) => s,
        Err(e) => {
            tracing::error!("failed to get article: {:?}", e);

            return render_not_found_error_tmpl(None);
        }
    };

    let front_matter_input =
        match fs::read_to_string(format!("./posts/{post_name}/article_frontmatter.toml")) {
            Ok(s) => s,
            Err(e) => {
                tracing::error!("failed to get front_matter: {:?}", e);

                return render_not_found_error_tmpl(None);
            }
        };

    let front_matter: FrontMatter = match toml::from_str(&front_matter_input) {
        Ok(fm) => fm,
        Err(e) => {
            tracing::error!("failed to convert front_matter: {:?}", e);

            return render_internal_error_tmpl(None);
        }
    };

    let parser = Parser::new_ext(&markdown_input, options);

    let mut html_output = String::new();
    html::push_html(&mut html_output, parser);

    let post_view = Post::new(crate::views::PostData {
        post: html_output,
        meta_data: front_matter,
    });

    render_template(post_view)
}
