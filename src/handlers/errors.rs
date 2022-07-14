use actix_web::{web, HttpResponse, Responder};
// use actix_web::{error, HttpResponse, http::{header::ContentType, StatusCode}};

// #[derive(Debug)]
// enum MyError {
//     InternalError,
//     BadClientData,
//     Timeout,
// }

// impl std::fmt::Display for MyError {
//     fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
//         todo!()
//     }
// }

// impl error::ResponseError for MyError {
//     fn error_response(self) -> HttpResponse {
//         HttpResponse::build(self.status_code())
//             .insert_header(ContentType::html())
//             .body(self.to_string())
//     }

//     fn status_code(&self) -> StatusCode {
//         match *self {
//             MyError::InternalError => StatusCode::INTERNAL_SERVER_ERROR,
//             MyError::BadClientData => StatusCode::BAD_REQUEST,
//             MyError::Timeout => StatusCode::GATEWAY_TIMEOUT,
//         }
//     }
// }

pub async fn not_found(tmpl: web::Data<tera::Tera>) -> impl Responder {
    let not_found_page = tmpl
        .render("not_found.html", &tera::Context::new())
        .unwrap();
    HttpResponse::InternalServerError()
        .content_type("text/html")
        .body(not_found_page)
}
