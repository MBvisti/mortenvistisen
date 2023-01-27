mod domain;
mod handlers;
mod service;

pub use domain::Email;
pub use handlers::{delete_subscriber, subscribe, verify_subscription};
