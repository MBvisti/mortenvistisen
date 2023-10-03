mod dashboard;
mod home;
mod post;
mod subscribe;

use actix_web::Responder;
pub use home::*;
pub use post::*;
pub use subscribe::*;

use crate::views::render_not_found_error_tmpl;

pub async fn not_found() -> impl Responder {
    render_not_found_error_tmpl(None)
}
