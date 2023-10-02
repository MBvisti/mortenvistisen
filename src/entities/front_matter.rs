use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct FrontMatter {
    pub title: String,
    pub file_name: String,
    pub description: String,
    pub posted: String,
    pub thumbnail: String,
    pub tags: Vec<String>,
    pub author: String,
    pub estimated_reading_time: u32,
    pub order: u32,
}
