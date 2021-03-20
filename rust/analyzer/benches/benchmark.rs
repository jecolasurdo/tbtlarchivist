use analyzer::engines::cosine_similarity::{self, Settings};
use analyzer::engines::Analyzer;
use criterion::{black_box, criterion_group, criterion_main, Criterion};
use std::fs::File;
use std::io::Read;

fn criterion_benchmark(c: &mut Criterion) {
    let sample_path = String::from(
        "/Users/Joe/Documents/code/tbtlarchivist/rust/analyzer/benches/drop_5000_samples.mp3",
    );
    let mut file = File::open(sample_path).unwrap();
    let mut data = Vec::new();
    file.read_to_end(&mut data).unwrap();

    // values in engine_settings are irrevent to benchmark
    let engine_settings = Settings {
        pass_one_sample_density: 1,
        pass_one_sample_size: 9,
        pass_one_threshold: 0.991,
        pass_two_sample_size: 50,
        pass_two_threshold: 0.99,
    };
    let engine = cosine_similarity::new(engine_settings);
    c.bench_function("mp3_to_raw", |b| {
        b.iter(|| engine.mp3_to_raw(black_box(&data)).unwrap())
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
