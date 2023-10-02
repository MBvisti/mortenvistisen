use crate::services::find_all_front_matter;
use actix_web::{get, Responder};

use crate::{
    // article::find_all_front_matter,
    views::{self, HomeIndexData},
};
use crate::views::{render_internal_error_tmpl, render_template};

#[tracing::instrument(name = "visit home page")]
#[get("/")]
pub async fn home_index() -> impl Responder {
    let mut front_matters = match find_all_front_matter() {
        Ok(fm) => fm,
        Err(e) => {
            println!("failed to find all frontmatters: {:?}", e);
            // tracing::error!("failed to find all frontmatters: {:?}", e);

            return render_internal_error_tmpl(None);
        }
    };

    front_matters.sort_by(|a, b| b.order.cmp(&a.order));
    
    let home_index_view = views::HomeIndex::new(HomeIndexData {
        posts: front_matters,
    });

    render_template(home_index_view)
}
