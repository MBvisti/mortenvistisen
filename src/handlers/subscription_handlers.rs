use actix_web::{post, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;

use crate::{domain::Email, repository};

#[derive(Debug, Serialize)]
struct SubscribeMetaData {
    has_error: bool,
    error_msg: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct SubscribeFormData {
    pub email: String,
    pub referer: String,
}

pub struct NewSubscriberPayload {
    pub email: Email,
    pub referer: String,
}

impl TryFrom<SubscribeFormData> for NewSubscriberPayload {
    type Error = String;

    fn try_from(value: SubscribeFormData) -> Result<Self, Self::Error> {
        let email = Email::parse(value.email)?;
        Ok(Self {
            email,
            referer: value.referer,
        })
    }
}

async fn render_subscribe_err(
    error_msg: String,
    mut context: tera::Context,
    tmpl: web::Data<tera::Tera>,
) -> HttpResponse {
    context.insert(
        "meta_data",
        &SubscribeMetaData {
            error_msg: Some(error_msg),
            has_error: true,
        },
    );

    match tmpl.render("_subscribe_response.html", &context) {
        Ok(s) => return HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            let error_page = tmpl.render("error_page.html", &context).unwrap();
            return HttpResponse::InternalServerError()
                .content_type("text/html")
                .body(error_page);
        }
    }
}
#[post("/subscribe")]
pub async fn subscribe(
    tmpl: web::Data<tera::Tera>,
    form: web::Form<SubscribeFormData>,
    pool: web::Data<PgPool>,
) -> impl Responder {
    let mut context = tera::Context::new();

    let payload: NewSubscriberPayload = match form.0.try_into() {
        Ok(p) => p,
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                tmpl,
            )
            .await;
        }
    };

    let email_exists = match repository::does_email_exists(&pool, &payload.email).await {
        Ok(exists) => exists,
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                tmpl,
            )
            .await;
        }
    };

    if email_exists {
        return render_subscribe_err("Email already registered".to_string(), context, tmpl).await;
    }

    if let Err(e) = repository::create_new_subscriber(&pool, &payload.email, &payload.referer).await
    {
        println!("{:?}", e);
        return render_subscribe_err(
            "I fucked up somehow, sorry. Please try again".to_string(),
            context,
            tmpl,
        )
        .await;
    };

    context.insert(
        "meta_data",
        &SubscribeMetaData {
            error_msg: None,
            has_error: false,
        },
    );

    match tmpl.render("_subscribe_response.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                tmpl,
            )
            .await;
        }
    }
}

#[post("/subscribe/verify")]
pub async fn verify_subscription(tmpl: web::Data<tera::Tera>) -> impl Responder {
    let context = tera::Context::new();

    match tmpl.render("post.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            let error_page = tmpl.render("error_page.html", &context).unwrap();
            HttpResponse::InternalServerError()
                .content_type("text/html")
                .body(error_page)
        }
    }
}
