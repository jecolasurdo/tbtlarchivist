#[global_allocator]
static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use analyzer::accessors::http;
use analyzer::engines::cosim_two_pass::{self, Settings};
use analyzer::managers;
use anyhow::Result;
use interop::BytesExt;
use log::LevelFilter;
use log4rs::append::file::FileAppender;
use log4rs::config::{Appender, Config, Root};
use log4rs::encode::pattern::PatternEncoder;
use protobuf::Message;
use std::io::{self, Read, Write};

fn main() -> Result<()> {
    let logfile = FileAppender::builder()
        .encoder(Box::new(PatternEncoder::new("{l} - {m}\n")))
        .build("analyzerd.log")?;

    let config = Config::builder()
        .appender(Appender::builder().build("logfile", Box::new(logfile)))
        .build(Root::builder().appender("logfile").build(LevelFilter::Info))?;

    log4rs::init_config(config)?;
    log::info!("Hello, World!");

    let engine_settings = Settings {
        target_sample_rate: 22_050,
        rms_window_size: 2756,
        pass_one_sample_size: 50,
        pass_one_threshold: 0.60,
        pass_two_sample_size: 500,
        pass_two_threshold: 0.9,
    };
    let analyzer_engine = cosim_two_pass::new(engine_settings);
    let uri_accessor = http::Accessor {};
    let mgr = managers::standard::new::<
        analyzer::engines::cosim_two_pass::Engine,
        analyzer::accessors::http::Accessor,
        analyzer::engines::cosim_two_pass::Error,
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
