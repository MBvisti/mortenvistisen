use super::{views::ViewData, View};

pub struct SubscribeResponse(ViewData);

impl View for SubscribeResponse {
    fn template_path(&self) -> &str {
        &self.0.path.as_ref()
    }

    fn get_context(&self) -> &tera::Context {
        &self.0.context
    }
}

pub struct SubscribeResponseData {
    pub has_error: bool,
    pub error_msg: String,
}
impl SubscribeResponse {
    pub fn new(payload: SubscribeResponseData) -> Self {
        let mut ctx = tera::Context::new();
        ctx.insert("has_error", &payload.has_error);
        ctx.insert("error_msg", &payload.error_msg);

        let view_data = ViewData::new(String::from("subscribe/_response.html"), ctx);
        return Self { 0: view_data };
    }
}

pub struct SubscribeVerify(ViewData);

impl View for SubscribeVerify {
    fn template_path(&self) -> &str {
        &self.0.path.as_ref()
    }

    fn get_context(&self) -> &tera::Context {
        &self.0.context
    }
}

pub struct SubscribeVerifyData {
    pub email_deleted: bool,
    pub has_error: bool,
    pub already_verified: bool,
    pub error_msg: Option<String>,
}
impl SubscribeVerify {
    pub fn new(payload: SubscribeVerifyData) -> Self {
        let mut ctx = tera::Context::new();
        ctx.insert("has_error", &payload.has_error);
        ctx.insert("error_msg", &payload.error_msg);
        ctx.insert("email_deleted", &payload.email_deleted);

        let view_data = ViewData::new(String::from("subscribe/email_confirm.html"), ctx);
        return Self { 0: view_data };
    }
}
