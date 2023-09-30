use tera::Context;

pub struct ViewData {
    pub path: String,
    pub context: Context,
}

impl ViewData {
    pub fn new(path: String, context: Context) -> Self {
        return ViewData { path, context };
    }
}

pub trait View {
    fn template_path(&self) -> &str;
    fn get_context(&self) -> &Context;
}
