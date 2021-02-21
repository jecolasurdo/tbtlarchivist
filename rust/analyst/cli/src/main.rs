#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

// use analyzer::{analyze, read_wav};
use anyhow::Result;
use std::io::{self, Read};

// const FILE_NAME_EPISODE: &str = "/Users/Joe/Documents/code/tbtldrops/audio/episodes/episode.wav";
// const FILE_NAME_DROP: &str = "/Users/Joe/Documents/code/tbtldrops/audio/drops/drop.wav";

fn main() -> Result<()> {
    // let episode = read_wav(FILE_NAME_EPISODE.to_owned())?;
    // println!("episode samples: {}", episode.len());

    // let drop = read_wav(FILE_NAME_DROP.to_owned())?;
    // println!("{}", drop.len());

    // let analyzer_options = analyzer::AnalyzerOptions {
    //     pass_one_sample_density: 1,
    //     pass_one_sample_size: 9,
    //     pass_one_threshold: 0.991,
    //     pass_two_sample_size: 50,
    //     pass_two_threshold: 0.99,
    // };

    // analyze(episode, drop, analyzer_options)?;
    
    let mut buffer = String::new();
    io::stdin().read_to_string(&mut buffer)?;
    print!("hello {}", buffer);
    Ok(())
}
