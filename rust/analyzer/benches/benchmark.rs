use analyzer::engines::cosine_similarity::{self, Settings};
use analyzer::engines::Analyzer;
use criterion::{black_box, criterion_group, criterion_main, Criterion};
use std::fs::File;
use std::io::Read;

fn benchmark_mp3_to_raw(c: &mut Criterion) {
    let sample_path = String::from(
        "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/125ms_constant_192kbps_joint_stereo.mp3",
    );
    let mut file = File::open(sample_path).unwrap();
    let mut data = Vec::new();
    file.read_to_end(&mut data).unwrap();

    // values in engine_settings are irrevent to benchmark
    let engine_settings = Settings {
        target_sample_rate: 22_050,
        rms_window_size: 2756,
        pass_one_sample_size: 50,
        pass_one_threshold: 0.60,
        pass_two_sample_size: 500,
        pass_two_threshold: 0.8,
    };
    let engine = cosine_similarity::new(engine_settings);
    c.bench_function("mp3_to_raw", |b| {
        b.iter(|| engine.mp3_to_raw(black_box(&data)))
    });
}

criterion_group!(benches, benchmark_mp3_to_raw);
criterion_main!(benches);
