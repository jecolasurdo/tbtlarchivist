#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use std::convert::TryFrom;
use anyhow::Result;
use std::io::{self, Read, Write};
use contracts::{CompletedResearchItem, PendingResearchItem};
use protobuf::Message;

fn main() -> Result<()> {
    // try to receive a message from stdin
    let mut buffer = vec![];
    io::stdin().read_to_end(&mut buffer)?;
    let pending_research_item: PendingResearchItem = Message::parse_from_bytes(&buffer)?;

    for n in 1..11 {
        // construct an outbound message using the inbound message's lease_id
        let mut completed_research_item= CompletedResearchItem::default();
        completed_research_item.lease_id = format!("{}_{}",pending_research_item.lease_id, n);

        // construct a message frame for the outbound message so the upstream service can parse it
        let completed_research_item_bytes = completed_research_item.write_to_bytes()?;
        let mut frame  = i32::try_from(completed_research_item_bytes.len())?.to_be_bytes().to_vec();
        frame.extend(&completed_research_item_bytes);

        // ship it
        io::stdout().write(&frame)?;
    }


    Ok(())
}

////////////////////////////////////////////////////////////////////////////////

// use analyzer::{analyze, read_wav};

// const FILE_NAME_EPISODE: &str = "/Users/Joe/Documents/code/tbtldrops/audio/episodes/episode.wav";
// const FILE_NAME_DROP: &str = "/Users/Joe/Documents/code/tbtldrops/audio/drops/drop.wav";

// fn main() -> Result<()> {
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
// }