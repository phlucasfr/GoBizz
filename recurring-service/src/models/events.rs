use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use sqlx::types::chrono::{DateTime, NaiveDate, Utc};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize, FromRow, Clone)]
pub struct Event {
    pub id: Uuid,
    pub customer_id: Uuid,
    pub name: String,
    #[serde(alias = "startDate")]
    pub start_date: NaiveDate,
    #[serde(alias = "intervalDays")]
    pub interval_days: i32,
    #[serde(alias = "stopAt")]
    pub stop_at: Option<NaiveDate>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(deny_unknown_fields)]
pub struct CreateEvent {
    pub name: String,
    #[serde(alias = "startDate")]
    pub start_date: NaiveDate,
    #[serde(alias = "intervalDays")]
    pub interval_days: i32,
    #[serde(alias = "stopAt")]
    pub stop_at: Option<NaiveDate>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(deny_unknown_fields)]
pub struct UpdateEvent {
    pub name: Option<String>,
    #[serde(alias = "startDate")]
    pub start_date: Option<sqlx::types::chrono::NaiveDate>,
    #[serde(alias = "intervalDays")]
    pub interval_days: Option<i32>,
    #[serde(default, with = "serde_with::rust::double_option", alias = "stopAt")]
    pub stop_at: Option<Option<sqlx::types::chrono::NaiveDate>>,
}
