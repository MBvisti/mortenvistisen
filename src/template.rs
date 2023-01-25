use tera::{Context, Tera};

lazy_static! {
    static ref TEMPLATES: Tera = {
        let mut tera = match Tera::new("templates/**/*.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {}", e);
                ::std::process::exit(1);
            }
        };
        tera.autoescape_on(vec![".html", ".sql"]);
        tera
    };
    static ref INTERNAL_SERVER_ERROR_TMPL: Tera = {
        let mut tera = match Tera::new("templates/error_templates/500.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {}", e);
                ::std::process::exit(1);
            }
        };
        tera.autoescape_on(vec![".html", ".sql"]);
        tera
    };
    static ref NOT_FOUND_ERROR_TMPL: Tera = {
        let mut tera = match Tera::new("templates/error_templates/404.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {}", e);
                ::std::process::exit(1);
            }
        };
        tera.autoescape_on(vec![".html", ".sql"]);
        tera
    };
}

pub fn render_template(template_name: &str, context: &Context) -> String {
    match TEMPLATES.render(template_name, context) {
        Ok(tmpl) => tmpl,
        Err(e) => {
            println!("error rendering template: {}", e);
            render_internal_error_tmpl(None)
        }
    }
}

pub fn render_internal_error_tmpl(provided_context: Option<&Context>) -> String {
    let context = tera::Context::new();
    if let Some(provided_context) = provided_context {
        return INTERNAL_SERVER_ERROR_TMPL
            .render("500.html", provided_context)
            .unwrap_or("something horrible happend if you see this".to_string())
    }
    INTERNAL_SERVER_ERROR_TMPL
        .render("500.html", &context)
        .unwrap_or("something horrible happend if you see this".to_string())
}

pub fn render_not_found_error_tmpl(provided_context: Option<&Context>) -> String {
    let context = tera::Context::new();
    if let Some(provided_context) = provided_context {
        return NOT_FOUND_ERROR_TMPL
            .render("404.html", provided_context)
            .unwrap_or("something horrible happend if you see this".to_string())
    }
    NOT_FOUND_ERROR_TMPL
        .render("404.html", &context)
        .unwrap_or("something horrible happend if you see this".to_string())
}
