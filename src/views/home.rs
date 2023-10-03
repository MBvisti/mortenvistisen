use crate::entities::FrontMatter;
use serde::{Deserialize, Serialize};

use super::{View, ViewData};

pub struct HomeIndex(ViewData);

impl View for HomeIndex {
    fn template_path(&self) -> &str {
        self.0.path.as_ref()
    }

    fn get_context(&self) -> &tera::Context {
        &self.0.context
    }
}

#[derive(Serialize, Debug, Deserialize)]
pub struct HomeIndexData {
    pub posts: Vec<FrontMatter>,
}

impl HomeIndex {
    pub fn new(data: HomeIndexData) -> Self {
        let mut ctx = tera::Context::new();
        ctx.insert("posts", &data.posts);

        let view_data = ViewData::new(String::from("home/index.html"), ctx);
        Self(view_data)
    }
}
