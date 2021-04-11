#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use analyzer::accessors::http;
use analyzer::engines::cosine_similarity::{self, Settings};
use analyzer::managers;
use anyhow::Result;
use interop::BytesExt;
use protobuf::Message;
use std::io::{self, Read, Write};

fn main() -> Result<()> {
    let engine_settings = Settings {
        target_sample_rate: 22_050,
        rms_window_size: 2756,
        pass_one_sample_size: 50,
        pass_one_threshold: 0.60,
        pass_two_sample_size: 500,
        pass_two_threshold: 0.9,
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

    let mut buffer = vec![];
    io::stdin().read_to_end(&mut buffer).unwrap();
    let pri = Message::parse_from_bytes(&buffer).unwrap();

    let ctx = cancel::Token::new();
    let rx = mgr.run(&ctx, &pri);

    while !ctx.is_canceled() {
        match rx.recv() {
            Ok(cri) => {
                let frame = cri?.write_to_bytes()?.to_frame();
                io::stdout().write(&frame)?;
            }
            Err(_) => {
                ctx.cancel();
            }
        }
    }

    Ok(())
}
