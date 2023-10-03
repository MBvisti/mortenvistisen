use actix_web::HttpResponse;
use tera::{Context, Tera};

use super::View;

const INTERNAL_SERVER_ERR_TMPL: &str = "errors/500.html";
const NOT_FOUND_ERR_TMPL: &str = "errors/404.html";

lazy_static! {
    static ref TEMPLATES: Tera = {
        let mut tera = match Tera::new("templates/**/*.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {e}");
                ::std::process::exit(1);
            }
        };
        tera.autoescape_on(vec![".html", ".sql"]);
        tera
    };
}

pub fn render_template(view: impl View) -> HttpResponse {
    match TEMPLATES.render(view.template_path(), view.get_context()) {
        Ok(tmpl) => HttpResponse::Ok().content_type("text/html").body(tmpl),
        Err(e) => {
            println!("Parsing error(s): {e}");

            render_internal_error_tmpl(None)
        }
    }
}

pub fn render_email_template(template_path: &str, ctx: &tera::Context) -> String {
    TEMPLATES.render(template_path, ctx).unwrap()
}

pub fn render_internal_error_tmpl(provided_context: Option<&Context>) -> HttpResponse {
    let mut context = &tera::Context::new();

    if let Some(provided_context) = provided_context {
        context = provided_context;
    }

    HttpResponse::InternalServerError()
        .content_type("text/html")
        .body(TEMPLATES.render(INTERNAL_SERVER_ERR_TMPL, context).unwrap())
}

pub fn render_not_found_error_tmpl(provided_context: Option<&Context>) -> HttpResponse {
    let mut context = &tera::Context::new();

    if let Some(provided_context) = provided_context {
        context = provided_context;
    }

    HttpResponse::Ok()
        .content_type("text/html")
        .body(TEMPLATES.render(NOT_FOUND_ERR_TMPL, context).unwrap())
}
