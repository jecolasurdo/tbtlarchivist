#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use std::convert::TryFrom;
// use anyhow::Result;
use std::io::{self, Read, Write};
use contracts::{PendingResearchItem, CompletedResearchItem};
use protobuf::Message;
use actix_web::client::Client;
use actix_web::http::StatusCode;

#[actix_web::main]
async fn main()  {
    frame_stderr_string(String::from("Starting..."));


    // try to receive a message from stdin
    let mut buffer = vec![];
    io::stdin().read_to_end(&mut buffer).unwrap();
    let mut pending_research_item = PendingResearchItem::default();
    match pending_research_item = Message::parse_from_bytes(&buffer) {
        Ok(pri) => pri,
        Err(err) => {
            frame_stderr_string(format!("{:?}", err));
            return ();
        };
    }

    frame_stderr_string(String::from("Received pending research item. Checking media type..."));
    
    let episode = pending_research_item.get_episode();
    if episode.get_media_type() != "mp3" {
        frame_stderr_string(String::from("the rust analyzer currently only supports mp3s"));
    }
    
    frame_stderr_string(String::from("Checked media type. Starting actix client and awaiting response..."));
   
    let client = Client::default();
    let response = client.get(episode.get_media_uri()).send().await;

    frame_stderr_string(String::from("Done awaiting response..."));

    match response {
        Ok(ref res) => {
            frame_stderr_string(format!("{:?}", res));
        },
        Err(err) => {
            frame_stderr_string(format!("{:?}", err));
            return ();
        }
    };

    if response.unwrap().status() != StatusCode::OK {
        frame_stderr_string(String::from("episode media URI did not return 200"));
        ()
    }

    // construct a message frame for the outbound message so the upstream service can parse it
    let mut completed_research_item = CompletedResearchItem::default();
    completed_research_item.lease_id = pending_research_item.lease_id;
    frame_stdout(completed_research_item.write_to_bytes().unwrap());

}

fn frame_stderr_string(s: String) {
    let f = frame_from_string(s);
    io::stderr().write(&f).unwrap();
}

fn frame_stdout(b: Vec<u8>) {
    let f = frame(b);
    io::stdout().write(&f).unwrap();
}

fn frame(b: Vec<u8>)-> Vec<u8> {
    let mut frame  = i32::try_from(b.len()).unwrap().to_be_bytes().to_vec();
    frame.extend(&b);
    frame
}

fn frame_from_string(s: String) -> Vec<u8> {
    return frame(s.as_bytes().to_vec())
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