use tera::{Context, Tera};

fn main() {
    let tera = match Tera::new("templates/**/*") {
        Ok(t) => t,
        Err(e) => {
            println!("Parsing error(s): {}", e);
            ::std::process::exit(1);
        }
    };
    let context = Context::new();
    let r = tera
        .render("index.html", &context)
        .expect("expected to be rendered");
    println!("{}", r)
}
