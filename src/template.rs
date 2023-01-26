use tera::{Context, Tera, Error};

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
}

pub fn render_template(template_name: &str, context: &Context) -> Result<String, Error> {
    match TEMPLATES.render(template_name, context) {
        Ok(tmpl) => Ok(tmpl),
        Err(e) => {
            println!("error rendering template: {}", e);
            Err(e)
        }
    }
}

pub fn render_internal_error_tmpl(provided_context: Option<&Context>) -> String {
    let context = tera::Context::new();
    if let Some(provided_context) = provided_context {
        return TEMPLATES
            .render("500.html", provided_context)
            .unwrap_or_else(|_| "something horrible happend if you see this".to_string());
    }
    TEMPLATES
        .render("500.html", &context)
        .unwrap_or_else(|_| "something horrible happend if you see this".to_string())
}

pub fn render_not_found_error_tmpl(provided_context: Option<&Context>) -> String {
    let context = tera::Context::new();
    if let Some(provided_context) = provided_context {
        return TEMPLATES
            .render("error_templates/404.html", provided_context)
            .unwrap_or_else(|_| "something horrible happend if you see this".to_string());
    }
    TEMPLATES
        .render("error_templates/404.html", &context)
        .unwrap_or_else(|_| "something horrible happend if you see this".to_string())
}
