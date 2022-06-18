use std::fs;

use actix_web::{get, web, HttpResponse, Responder};
use ignore::WalkBuilder;
use serde::{Deserialize, Serialize};
// use pulldown_cmark::{html, Options, Parser};

#[derive(Serialize, Deserialize, Debug)]
struct FrontMatter {
    title: String,
    description: String,
    posted: String,
    thumbnail: String,
    tags: Vec<String>,
    author: String,
    estimated_reading_time: u32,
}

#[get("/")]
pub async fn index(tmpl: web::Data<tera::Tera>) -> impl Responder {
    let mut context = tera::Context::new();
    // Set up options and parser. Strikethroughs are not part of the CommonMark standard
    // and we therefore must enable it explicitly.
    // let mut options = Options::empty();
    // // options.insert(Options::ENABLE_STRIKETHROUGH);
    // let markdown_input = match fs::read_to_string("./posts/practical-approach-to-structure-go-apps.md") {
    //     Ok(s) => s,
    //     Err(_) => "hello".to_string(),
    // };
    //
    let mut t = ignore::types::TypesBuilder::new();
    t.add_defaults();
    t.select("toml");
    let mut front_matters = Vec::new();
    // t.build().unwrap();
    // let parser = Parser::new_ext(&markdown_input, options);
    for res in WalkBuilder::new("./posts")
        .types(t.build().unwrap())
        .build()
    {
        match res {
            Ok(entry) => {
                if entry.path().is_file() {
                    let file_content = match fs::read_to_string(entry.path()) {
                        Ok(s) => s,
                        Err(e) => {
                            println!("{:?}", e.to_string());
                            e.to_string()
                        }
                    };

                    let front_matter: FrontMatter = match toml::from_str(&file_content) {
                        Ok(s) => s,
                        Err(e) => {
                            panic!("{:?}", e.to_string());
                        }
                    };
                    front_matters.push(front_matter);
                }
            }
            Err(err) => println!("ERROR: {}", err),
        }
    }
    // for result in Walk::new("./posts") {
    //     // Each item yielded by the iterator is either a directory entry or an
    //     // error, so either print the path or the error.
    //     match result {
    //         Ok(entry) => println!("{}", entry.path().display()),
    //         Err(err) => println!("ERROR: {}", err),
    //     }
    // }
    // // Write to String buffer.
    // let html_output = String::new();
    // html::push_html(&mut html_output, parser);

    // println!("{:?}", front_matters);
    context.insert("posts", &front_matters);
    // let c = &Context::from_serialize(front_matters).unwrap();
    match tmpl.render("home.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            HttpResponse::InternalServerError().content_type("text/html").body("")
        }
    }
}
