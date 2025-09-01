use serde::{Deserialize, Serialize};
use sqlx::types::chrono::NaiveDate;
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct OccurrenceDTO {
    pub event_id: Uuid,
    pub name: String,
    pub date: NaiveDate,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct OccurrencePeriod {
    pub start: NaiveDate,
    pub end: NaiveDate,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct OccurrencesResponse {
    pub period: OccurrencePeriod,
    pub occurrences: Vec<OccurrenceDTO>,
}
