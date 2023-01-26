use actix_web::{get, post, web, HttpResponse, Responder};
use rand::distributions::Alphanumeric;
use rand::{thread_rng, Rng};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use uuid::Uuid;

use crate::email_client::EmailClient;
use crate::template::{render_internal_error_tmpl, render_template};
use crate::{domain::Email, repository};

#[derive(Debug, Serialize)]
struct SubscribeMetaData {
    has_error: bool,
    error_msg: Option<String>,
}

#[derive(Debug, Serialize)]
struct ConfirmSubscribeMetaData {
    has_error: bool,
    already_verified: bool,
    error_msg: Option<String>,
}

#[derive(Debug, Serialize)]
struct ConfirmSubEmail {
    app_base_url: String,
    token: String,
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
    template_name: &str,
) -> HttpResponse {
    context.insert(
        "meta_data",
        &SubscribeMetaData {
            error_msg: Some(error_msg),
            has_error: true,
        },
    );

    let tmpl = match render_template(template_name, &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}

fn generate_subcription_token() -> String {
    let mut rng = thread_rng();
    std::iter::repeat_with(|| rng.sample(Alphanumeric))
        .map(char::from)
        .take(25)
        .collect()
}

#[post("/subscribe")]
pub async fn subscribe(
    form: web::Form<SubscribeFormData>,
    pool: web::Data<PgPool>,
    email_client: web::Data<EmailClient>,
) -> impl Responder {
    let mut context = tera::Context::new();
    let mut email_context = tera::Context::new();

    let payload: NewSubscriberPayload = match form.0.try_into() {
        Ok(p) => p,
        Err(e) => {
            println!("$$$$$$$$$$$$$$$$$$$$$$$$$$$$ YOYOY: {:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "_subscribe_response.html",
            )
            .await;
        }
    };

    let email_exists = match repository::does_email_exists(&pool, &payload.email).await {
        Ok(exists) => exists,
        Err(e) => {
            println!("$$$$$$$$$$$$$$$$$$$$$$$$$$$$ YOYOY: {:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "_subscribe_response.html",
            )
            .await;
        }
    };
    if email_exists {
        return render_subscribe_err(
            "Email already registered".to_string(),
            context,
            "_subscribe_response.html",
        )
        .await;
    }
    let new_sub_id = Uuid::new_v4();
    if let Err(e) =
        repository::create_new_subscriber(&pool, &payload.email, &payload.referer, new_sub_id).await
    {
        println!("{:?}", e);
        return render_subscribe_err(
            "I fucked up somehow, sorry. Please try again".to_string(),
            context,
            "_subscribe_response.html",
        )
        .await;
    };
    let new_sub_token = generate_subcription_token();
    if let Err(e) = repository::create_new_subscriber_token(&pool, new_sub_id, &new_sub_token).await
    {
        println!("{:?}", e);
        return render_subscribe_err(
            "I fucked up somehow, sorry. Please try again".to_string(),
            context,
            "_subscribe_response.html",
        )
        .await;
    };

    email_context.insert(
        "meta_data",
        &ConfirmSubEmail {
            app_base_url: email_client.app_base_url.clone(),
            token: new_sub_token,
        },
    );
    let html_content = match render_template("confirm_sub_email.html", &email_context) {
        Ok(t) => t,
        Err(_) => return HttpResponse::InternalServerError().content_type("text/html").body(render_internal_error_tmpl(None))
    };
    match email_client
        .send_email(
            payload.email,
            "Thanks for subscribing to my newsletter!",
            &html_content,
        )
        .await
    {
        Ok(_) => (),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "_subscribe_response.html",
            )
            .await;
        }
    };

    context.insert(
        "meta_data",
        &SubscribeMetaData {
            error_msg: None,
            has_error: false,
        },
    );

    let tmpl = match render_template("_subscribe_response.html", &context) {
        Ok(t) => t,
        Err(_) => render_internal_error_tmpl(None),
    };

    HttpResponse::Ok().content_type("text/html").body(tmpl)
}

#[derive(serde::Deserialize)]
pub struct Parameters {
    token: String,
}

#[get("/subscribe/verify")]
pub async fn verify_subscription(
    tmpl: web::Data<tera::Tera>,
    pool: web::Data<PgPool>,
    params: web::Query<Parameters>,
) -> impl Responder {
    let mut context = tera::Context::new();

    let id = match repository::get_subscriber_id_from_token(&pool, &params.token).await {
        Ok(id) => id,
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "confirm_subscription.html",
            )
            .await;
        }
    };

    let is_verified = match repository::is_user_verified(&pool, id).await {
        Ok(is_verified) => is_verified,
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "confirm_subscription.html",
            )
            .await;
        }
    };

    if is_verified {
        match tmpl.render("confirm_subscription.html", &context) {
            Ok(s) => {
                context.insert(
                    "meta_data",
                    &ConfirmSubscribeMetaData {
                        already_verified: true,
                        has_error: false,
                        error_msg: None,
                    },
                );
                return HttpResponse::Ok().content_type("text/html").body(s);
            }
            Err(e) => {
                println!("{:?}", e);
                return render_subscribe_err(
                    "I fucked up somehow, sorry. Please try again".to_string(),
                    context,
                    "confirm_subscription.html",
                )
                .await;
            }
        }
    }

    match repository::update_email_to_verified(&pool, id).await {
        Ok(_) => (),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "confirm_subscription.html",
            )
            .await;
        }
    };

    match tmpl.render("confirm_subscription.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "confirm_subscription.html",
            )
            .await;
        }
    }
}

#[get("/subscribe/delete")]
pub async fn delete_subscriber(
    tmpl: web::Data<tera::Tera>,
    pool: web::Data<PgPool>,
    params: web::Query<Parameters>,
) -> impl Responder {
    let context = tera::Context::new();

    let id = match repository::get_subscriber_id_from_token(&pool, &params.token).await {
        Ok(id) => id,
        Err(e) => {
            println!("{:?}", e);
            match e {
                // TODO: add message that it has already been deleted
                sqlx::Error::RowNotFound => {
                    match tmpl.render("delete_subscription.html", &context) {
                        Ok(s) => return HttpResponse::Ok().content_type("text/html").body(s),
                        Err(e) => {
                            println!("{:?}", e);
                            return render_subscribe_err(
                                "I fucked up somehow, sorry. Please try again".to_string(),
                                context,
                                "confirm_subscription.html",
                            )
                            .await;
                        }
                    }
                }
                _ => {
                    return render_subscribe_err(
                        "I fucked up somehow, sorry. Please try again".to_string(),
                        context,
                        "delete_subscription.html",
                    )
                    .await
                }
            }
        }
    };

    match repository::delete_subscriber_token(&pool, &params.token).await {
        Ok(_) => (),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "delete_subscription.html",
            )
            .await;
        }
    };

    match repository::delete_subscriber(&pool, id).await {
        Ok(_) => (),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "delete_subscription.html",
            )
            .await;
        }
    };

    match tmpl.render("delete_subscription.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            return render_subscribe_err(
                "I fucked up somehow, sorry. Please try again".to_string(),
                context,
                "confirm_subscription.html",
            )
            .await;
        }
    }
}
