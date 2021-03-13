#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use analyzer::accessors::http;
use analyzer::engines::cosine_similarity::{self, Settings};
use analyzer::managers;
use anyhow::Result;

fn main() -> Result<()> {
    let engine_settings = Settings {
        pass_one_sample_density: 0,
        pass_one_sample_size: 0,
        pass_one_threshold: 0.0,
        pass_two_sample_size: 0,
        pass_two_threshold: 0.0,
    };
    let analyzer_engine = cosine_similarity::new(engine_settings);
    let uri_accessor = http::Accessor {};
    let mgr = managers::standard::new::<
        analyzer::engines::cosine_similarity::Engine,
        analyzer::accessors::http::Accessor,
        analyzer::engines::cosine_similarity::Error,
        analyzer::accessors::http::Error,
        anyhow::Error,
    >(analyzer_engine, uri_accessor);

    let ctx = cancel::Token::new();
    let pri = contracts::PendingResearchItem::default();

    let _rx = mgr.run(&ctx, &pri);

    Ok(())
}
