use std::{fs, io::Error};

use actix_web::{get, HttpResponse, Responder};
use ignore::WalkBuilder;
use serde::{Deserialize, Serialize};

use crate::template::{render_internal_error_tmpl, render_template};

#[derive(Serialize, Deserialize, Debug)]
pub struct FrontMatter {
    title: String,
    file_name: String,
    description: String,
    posted: String,
    thumbnail: String,
    tags: Vec<String>,
    author: String,
    estimated_reading_time: u32,
    order: u32,
}

fn find_all_front_matter() -> Result<Vec<FrontMatter>, Error> {
    let mut t = ignore::types::TypesBuilder::new();
    t.add_defaults();
    let toml = match t.select("toml").build() {
        Ok(t) => t,
        Err(e) => {
            println!("{:}", e); // just print the error for now
            return Err(Error::new(
                std::io::ErrorKind::Other,
                "could not build toml file type matcher",
            ));
        }
    };

    let file_walker = WalkBuilder::new("./posts").types(toml).build();

    let mut front_matters = Vec::new();
    for front_matter in file_walker {
        match front_matter {
            Ok(fm) => {
                if fm.path().is_file() {
                    let fm_content = fs::read_to_string(fm.path())?;
                    let front_matter: FrontMatter = toml::from_str(&fm_content)?;

                    front_matters.push(front_matter);
                }
            }
            Err(e) => {
                println!("{:}", e); // just print the error for now
                return Err(Error::new(
                    std::io::ErrorKind::NotFound,
                    "could not locate frontmatter",
                ));
            }
        }
    }

    Ok(front_matters)
}

#[get("/")]
pub async fn index() -> impl Responder {
    let mut context = tera::Context::new();

    let mut front_matters = match find_all_front_matter() {
        Ok(fm) => fm,
        Err(e) => {
            println!("{:?}", e);

            return HttpResponse::InternalServerError()
                .content_type("text/html")
                .body(render_internal_error_tmpl(None));
        }
    };

    front_matters.sort_by(|a, b| b.order.cmp(&a.order));
    context.insert("posts", &front_matters);

    let tmpl = match render_template("home.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}
