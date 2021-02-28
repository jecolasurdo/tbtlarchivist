#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use std::convert::TryFrom;
// use anyhow::Result;
use std::io::{self, Read, Write};
use contracts::{PendingResearchItem, CompletedResearchItem};
use protobuf::Message;
use actix_web::client::Client;
// use actix_web::http::StatusCode;

#[actix_web::main]
async fn main()  {
    // try to receive a message from stdin
    let mut buffer = vec![];
    io::stdin().read_to_end(&mut buffer).unwrap();
    let pending_research_item: PendingResearchItem = Message::parse_from_bytes(&buffer).unwrap();
    
    let episode = pending_research_item.get_episode();
    if episode.get_media_type() != "mp3" {
        let error = construct_frame("the rust analyzer currently only supports mp3s".as_bytes().to_vec());
        io::stderr().write(&error).unwrap();
    }
    
    let client = Client::default();
    let response = client.get(episode.get_media_uri()).send().await;
    
    response.and_then(|response| {
        let notreallyanerror = construct_frame(format!("{:?}", response).as_bytes().to_vec());
        io::stderr().write(&notreallyanerror).unwrap();
        Ok(())
    }).unwrap();


    // if response.status() != StatusCode::OK {
    //     panic!("episode media URI did not return 200")
    // }

    // construct a message frame for the outbound message so the upstream service can parse it
    let completed_research_item = CompletedResearchItem::default();
    let frame = construct_frame(completed_research_item.write_to_bytes().unwrap());

    // ship it
    io::stdout().write(&frame).unwrap();
}

fn construct_frame(b: Vec<u8>)-> Vec<u8> {
    let mut frame  = i32::try_from(b.len()).unwrap().to_be_bytes().to_vec();
    frame.extend(&b);
    frame
}

////////////////////////////////////////////////////////////////////////////////

// use analyzer::{analyze, read_wav};

// const FILE_NAME_EPISODE: &str = "/Users/Joe/Documents/code/tbtldrops/audio/episodes/episode.wav";
// const FILE_NAME_DROP: &str = "/Users/Joe/Documents/code/tbtldrops/audio/drops/drop.wav";

// fn main() -> Result<()>{
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