#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use analyzer::accessors::file;
use analyzer::engines::cosine_similarity::{self, Settings};
use analyzer::managers;
use anyhow::Result;
use protobuf::RepeatedField;

fn main() -> Result<()> {
    let engine_settings = Settings {
        pass_one_sample_size: 50,
        pass_one_threshold: 0.60,
        pass_two_sample_size: 500,
        pass_two_threshold: 0.9,
    };
    let analyzer_engine = cosine_similarity::new(engine_settings);
    let uri_accessor = file::Accessor {};
    let mgr = managers::standard::new::<
        analyzer::engines::cosine_similarity::Engine,
        analyzer::accessors::file::Accessor,
        analyzer::engines::cosine_similarity::Error,
        analyzer::accessors::file::Error,
        anyhow::Error,
    >(analyzer_engine, uri_accessor);

    let ctx = cancel::Token::new();
    let mut episode = contracts::EpisodeInfo::default();
    episode.set_media_uri(String::from(
        "/Users/Joe/Documents/code/tbtlarchivist/rust/audio/episodes/episode.mp3",
    ));

    let mut clip = contracts::ClipInfo::default();
    clip.set_media_uri(String::from(
        "/Users/Joe/Documents/code/tbtlarchivist/rust/audio/drops/drop.mp3",
    ));
    let mut pri = contracts::PendingResearchItem::default();
    pri.set_lease_id(String::from("test_lease_id"));
    pri.set_episode(episode);
    pri.set_clips(RepeatedField::from_vec(vec![clip]));

    let rx = mgr.run(&ctx, &pri);

    while !ctx.is_canceled() {
        match rx.recv() {
            Ok(cri) => println!("{:?}", cri.unwrap()),
            Err(e) => {
                println!("Error: {:?}", e);
                ctx.cancel();
            }
        }
    }

    Ok(())
}
