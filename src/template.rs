use tera::{Context, Tera};

const INTERNAL_SERVER_ERR_TMPL: String = "templates/errors/500.html".to_string();
const NOT_FOUND_ERR_TMPL: String = "templates/errors/404.html".to_string();

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

pub fn render_template(template_name: &str, context: &Context) -> String {
    match TEMPLATES.render(template_name, context) {
        Ok(tmpl) => tmpl,
        Err(e) => {
            println!("Parsing error(s): {e}");
            render_internal_error_tmpl(None)
        }
    }
}

pub fn render_internal_error_tmpl(provided_context: Option<&Context>) -> String {
    let context = &tera::Context::new();

    if let Some(provided_context) = provided_context {
        context = provided_context;
    }

    Tera::one_off(&INTERNAL_SERVER_ERR_TMPL, context, true).unwrap()
}

pub fn render_not_found_error_tmpl(provided_context: Option<&Context>) -> String {
    let context = &tera::Context::new();

    if let Some(provided_context) = provided_context {
        context = provided_context;
    }

    Tera::one_off(&NOT_FOUND_ERR_TMPL, context, true).unwrap()
}
