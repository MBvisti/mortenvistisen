use serde::{Deserialize, Serialize};

use crate::entities::FrontMatter;

use super::{views::ViewData, View};

pub struct Post(ViewData);

impl View for Post {
    fn template_path(&self) -> &str {
        &self.0.path.as_ref()
    }

    fn get_context(&self) -> &tera::Context {
        &self.0.context
    }
}

#[derive(Serialize, Debug, Deserialize)]
pub struct PostData {
    pub post: String,
    pub meta_data: FrontMatter,
}

impl Post {
    pub fn new(data: PostData) -> Self {
        let mut ctx = tera::Context::new();
        ctx.insert("post", &data.post);
        ctx.insert("meta_data", &data.meta_data);

        let view_data = ViewData::new(String::from("post/post.html"), ctx);
        return Self { 0: view_data };
    }
}
