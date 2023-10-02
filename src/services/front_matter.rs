use crate::entities::FrontMatter;
use ignore::WalkBuilder;
use std::{fs, io::Error};

pub fn find_all_front_matter() -> Result<Vec<FrontMatter>, Error> {
    let mut t = ignore::types::TypesBuilder::new();
    t.add_defaults();
    let toml = match t.select("toml").build() {
        Ok(t) => t,
        Err(e) => {
            println!("{e:}"); // just print the error for now
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
                println!("{e:}"); // just print the error for now
                return Err(Error::new(
                    std::io::ErrorKind::NotFound,
                    "could not locate frontmatter",
                ));
            }
        }
    }

    Ok(front_matters)
}
